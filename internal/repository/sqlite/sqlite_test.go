//+build manual

package sqlite_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/mustafadubul/product/internal/domain"
	"github.com/mustafadubul/product/internal/repository/sqlite"
	"github.com/stretchr/testify/assert"
)

type testDB struct {
	db *gorm.DB
}

func StartTestDB(t *testing.T) *sqlite.DB {
	db, err := sqlite.Open(true, "")
	assert.Nil(t, err)

	db.AutoMigrate(&domain.Product{})

	return sqlite.New(db)
}

func TestCreateEntry(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	p, err := db.Create(&domain.Product{
		ID:       1,
		ItemName: "camera",
		Lat:      99,
		Lng:      66,
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, p.ID)
}

func TestGetProduct(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	expectedProduct := &domain.Product{
		ItemName: "camera",
		Lat:      99,
		Lng:      66,
	}

	p, err := db.Create(expectedProduct)
	assert.Nil(t, err)

	product, err := db.Get(p.ID)
	assert.Nil(t, err)

	assert.Equal(t, expectedProduct, product)
}

func TestDeleteProduct(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	product := &domain.Product{
		ItemName: "camera",
		Lat:      99,
		Lng:      66,
	}

	p, err := db.Create(product)
	assert.Nil(t, err)

	err = db.Delete(p.ID)
	assert.Nil(t, err)

	_, err = db.Get(p.ID)
	assert.NotNil(t, err)
}

func TestUpdateProduct(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	oldProduct := &domain.Product{
		ID:       1,
		ItemName: "camera",
		Lat:      99,
		Lng:      66,
	}
	p, err := db.Create(oldProduct)
	assert.Nil(t, err)

	updatedProduct := &domain.Product{
		ID:       p.ID,
		ItemName: "Canon",
		Lat:      11,
		Lng:      22,
	}

	newP, err := db.Update(updatedProduct)
	assert.Nil(t, err)

	assert.Equal(t, updatedProduct, newP)
}

func TestSearchBetween(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	products := []*domain.Product{
		{
			ItemName: "camera london",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "camera london",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "camera london",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "camera paris",
			Lat:      48.864716,
			Lng:      2.349014,
		},
	}

	for _, p := range products {
		_, err := db.Create(p)
		assert.Nil(t, err)
	}

	points := []domain.Point{
		{
			51.509909966080286, -0.11809199999997985,
		}, {
			51.50984485186806, -0.11801975142352036,
		}, {
			51.50983808958223, -0.11809199999997985,
		}, {
			51.509909263797134, -0.1181642486786242,
		},
	}

	between := db.Between(points)

	p, err := db.Search(between)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(p))
}

func TestSearchLike(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	products := []*domain.Product{
		{
			ItemName: "Canon 50mm f/1.2 Prime Lens",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "Canon 6D +24-70mm F4 L/35mm + Microphone + LED Lights",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "Go Pro Hero - Full HD",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "Canon 5D Mii Shooting Kit and 28mm, 50mm and 105mm Lenses",
			Lat:      48.864716,
			Lng:      2.349014,
		},
	}

	for _, p := range products {
		_, err := db.Create(p)
		assert.Nil(t, err)
	}

	like := db.Like("Canon")

	p, err := db.Search(like)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(p))
}

func TestSearchBetweenLike(t *testing.T) {
	db := StartTestDB(t)
	defer db.Close()

	products := []*domain.Product{
		{
			ItemName: "Canon 50mm f/1.2 Prime Lens",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "Canon 6D +24-70mm F4 L/35mm + Microphone + LED Lights",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "Go Pro Hero - Full HD",
			Lat:      51.509865,
			Lng:      -0.118092,
		},
		{
			ItemName: "Go Pro Hero - Full HD",
			Lat:      48.864716,
			Lng:      2.349014,
		},
		{
			ItemName: "Canon 5D Mii Shooting Kit and 28mm, 50mm and 105mm Lenses",
			Lat:      48.864716,
			Lng:      2.349014,
		},
		{
			ItemName: "Canon 5D Mii Shooting Kit and 28mm, 50mm and 105mm Lenses",
			Lat:      48.864716,
			Lng:      2.349014,
		},
		{
			ItemName: "Canon 5D Mii Shooting Kit and 28mm, 50mm and 105mm Lenses",
			Lat:      48.864716,
			Lng:      2.349014,
		},
	}

	for _, p := range products {
		_, err := db.Create(p)
		assert.Nil(t, err)
	}

	points := []domain.Point{
		{
			51.509909966080286, -0.11809199999997985,
		}, {
			51.50984485186806, -0.11801975142352036,
		}, {
			51.50983808958223, -0.11809199999997985,
		}, {
			51.509909263797134, -0.1181642486786242,
		},
	}

	between := db.Between(points)

	like := db.Like("Go Pro Hero")

	p, err := db.Search(between, like)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(p))
}
