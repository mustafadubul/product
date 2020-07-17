package repository

import (
	"errors"

	"github.com/mustafadubul/product/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
	ErrFatal    = errors.New("fatal error")
)

// mockgen -source=repository.go -package=mocks -mock_names Product=MockRepoProduct -destination=../../mocks/mocks_repo_product.go Product
type Product interface {
	Search(functions ...func() error) ([]domain.Product, error)
	Like(term string) func() error
	Between(points []domain.Point) func() error

	Create(p *domain.Product) (*domain.Product, error)
	Get(id uint64) (*domain.Product, error)

	Update(p *domain.Product) (*domain.Product, error)
	Delete(id uint64) error
}
