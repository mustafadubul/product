package sqlite

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mustafadubul/product/internal/domain"
	"github.com/mustafadubul/product/internal/repository"
	"github.com/rs/zerolog"
)

type DB struct {
	db     *gorm.DB
	tx     *gorm.DB
	logger *zerolog.Logger
}

func Open(inMemory bool, fileName string) (*gorm.DB, error) {
	var file string
	if inMemory == true {
		file = "file::memory:?cache=shared"
	} else {
		file = fileName
	}

	return gorm.Open("sqlite3", file)
}

func New(db *gorm.DB) *DB {
	componentLogger := zerolog.New(os.Stdout).With().Str("component", "repository").Logger()
	return &DB{
		db:     db,
		logger: &componentLogger,
	}
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Migrate() {
	d.db.AutoMigrate(&domain.Product{})
}

func (d *DB) Search(functions ...func() error) ([]domain.Product, error) {
	var products []domain.Product

	d.tx = d.db
	for _, f := range functions {
		if err := f(); err != nil {
			return nil, err
		}
	}

	if err := d.tx.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (d *DB) Like(term string) func() error {
	return func() error {
		d.tx = d.tx.Where("item_name LIKE ? ", "%"+term+"%")
		return d.tx.Error
	}
}

func (d *DB) Between(points []domain.Point) func() error {
	return func() error {
		d.tx = d.tx.Where("lat > ? AND lat < ? AND lng < ? AND lng > ?",
			points[2].X, points[0].X, points[1].Y, points[3].Y)
		return d.tx.Error
	}
}

func (d *DB) Create(p *domain.Product) (*domain.Product, error) {
	if err := d.db.Create(p).Error; err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", repository.ErrFatal)
	}
	return p, nil
}

func (d *DB) Get(id uint64) (*domain.Product, error) {
	var product domain.Product

	err := d.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found product: %w", repository.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get product: %w", repository.ErrFatal)
	}

	return &product, nil
}

func (d *DB) Update(p *domain.Product) (*domain.Product, error) {
	if err := d.db.Model(&domain.Product{ID: p.ID}).Update(p).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %w", repository.ErrFatal)
	}
	return p, nil
}

func (d *DB) Delete(id uint64) error {
	if err := d.db.Delete(&domain.Product{ID: id}).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", repository.ErrFatal)
	}
	return nil
}
