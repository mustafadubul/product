package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mustafadubul/product/internal/domain"
	"github.com/mustafadubul/product/internal/service"
	"github.com/mustafadubul/product/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type Service struct {
	*service.Service
	ctrl   *gomock.Controller
	cancel context.CancelFunc

	mockProductRepo *mocks.MockRepoProduct
}

func CreateService(t *testing.T) *Service {
	ctrl := gomock.NewController(t)
	mockProductRepo := mocks.NewMockRepoProduct(ctrl)

	_, cancel := context.WithCancel(context.Background())
	l := zerolog.Nop()
	return &Service{
		Service:         service.New(&l, mockProductRepo),
		mockProductRepo: mockProductRepo,
		cancel:          cancel,
	}
}

func (s *Service) Finish() {
	s.ctrl.Finish()
	s.cancel()
}

func TestService(t *testing.T) {
	t.Run("query products", testSearch_QueryProducts)
	t.Run("query products with term", testSearch_QueryProductsWithterm)
	t.Run("get single products", testGetProduct)
	t.Run("update single products", testUpdateProduct)
	t.Run("delete single products", testDeleteProduct)
	t.Run("createw single products", testCreateProduct)
}

func testSearch_QueryProducts(t *testing.T) {
	s := CreateService(t)

	query := &domain.Query{
		Lat:    51.509865,
		Lng:    -0.118092,
		Radius: 5,
	}

	points := []domain.Point{
		{51.509909966080286, -0.11809199999997985},
		{51.50984485186806, -0.11801975142352036},
		{51.50983808958223, -0.11809199999997985},
		{51.509909263797134, -0.1181642486786242},
	}

	s.mockProductRepo.EXPECT().Between(points)

	products := []domain.Product{
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
	}

	s.mockProductRepo.EXPECT().Search(gomock.Any()).Return(products, nil)

	products, err := s.Search(context.Background(), query)
	assert.Nil(t, err)
	assert.NotNil(t, products)
}

func testSearch_QueryProductsWithterm(t *testing.T) {
	s := CreateService(t)

	query := &domain.Query{
		Lat:    51.509865,
		Lng:    -0.118092,
		Radius: 5,
		Term:   "canon",
	}

	points := []domain.Point{
		{51.509909966080286, -0.11809199999997985},
		{51.50984485186806, -0.11801975142352036},
		{51.50983808958223, -0.11809199999997985},
		{51.509909263797134, -0.1181642486786242},
	}

	s.mockProductRepo.EXPECT().Between(points)
	s.mockProductRepo.EXPECT().Like("canon")
	products := []domain.Product{
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
	}

	s.mockProductRepo.EXPECT().Search(gomock.Any()).Return(products, nil)

	products, err := s.Search(context.Background(), query)
	assert.Nil(t, err)
	assert.NotNil(t, products)
}

func testUpdateProduct(t *testing.T) {
	s := CreateService(t)
	p := &domain.Product{
		ID:       1234,
		ItemName: "canon",
	}
	s.mockProductRepo.EXPECT().Create(p).Return(p, nil)
}

func testDeleteProduct(t *testing.T) {
	s := CreateService(t)
	p := &domain.Product{
		ID: 1234,
	}
	s.mockProductRepo.EXPECT().Delete(p.ID).Return(nil)
}

func testGetProduct(t *testing.T) {
	s := CreateService(t)
	p := &domain.Product{
		ID:       1234,
		ItemName: "canon",
		Lat:      1234,
		Lng:      2123,
	}

	s.mockProductRepo.EXPECT().Get(p.ID).Return(p, nil)
}

func testCreateProduct(t *testing.T) {
	s := CreateService(t)
	p := &domain.Product{
		ItemName: "canon",
		Lat:      1234,
		Lng:      2123,
	}
	s.mockProductRepo.EXPECT().Create(p).Return(p, nil)
}
