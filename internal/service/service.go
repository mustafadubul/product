package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mustafadubul/product/internal/domain"

	"github.com/mustafadubul/product/internal/repository"
	"github.com/mustafadubul/product/pkg/geo"

	"github.com/rs/zerolog"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrRequestFailed = errors.New("request failed")
	ErrInputInvalid  = errors.New("input invalid")
)

type Service struct {
	ctx    context.Context
	logger *zerolog.Logger

	products repository.Product
}

func New(l *zerolog.Logger, productRepo repository.Product) *Service {
	componentLogger := l.With().Str("component", "service").Logger()
	return &Service{
		logger:   &componentLogger,
		products: productRepo,
	}
}

func (s *Service) Search(ctx context.Context, q *domain.Query) ([]domain.Product, error) {
	l := s.logger.With().Str("service", "Search").Logger()

	queryFunctions := []func() error{}

	points := geo.BoundingBox(q.Lat, q.Lng, q.Radius)

	between := s.products.Between(points)
	queryFunctions = append(queryFunctions, between)

	if q.Term != "" {
		like := s.products.Like(q.Term)
		queryFunctions = append(queryFunctions, like)
	}

	products, err := s.products.Search(queryFunctions...)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			l.Error().Err(err).Msg("failed to search products")
			return nil, fmt.Errorf("failed to search products: %w", ErrRequestFailed)
		}
		return nil, fmt.Errorf("product not found: %w", ErrNotFound)
	}

	return products, nil
}

func (s *Service) Update(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	l := s.logger.With().Str("service", "Update").Logger()

	p, err := s.products.Update(p)
	if err != nil {
		l.Error().Err(err).Msg("failed to update products")
		return nil, fmt.Errorf("failed to update products: %w", ErrRequestFailed)
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	l := s.logger.With().Str("service", "Delete").Logger()

	err := s.products.Delete(id)
	if err != nil {
		l.Error().Err(err).Msg("failed to delete products")
		return fmt.Errorf("failed to delete products: %w", ErrRequestFailed)
	}
	return nil
}

func (s *Service) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	l := s.logger.With().Str("service", "Create").Logger()

	p, err := s.products.Create(p)
	if err != nil {
		l.Error().Err(err).Msg("failed to create products")
		return nil, fmt.Errorf("failed to create products: %w", ErrRequestFailed)
	}
	return p, nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*domain.Product, error) {
	l := s.logger.With().Str("service", "Get").Logger()

	product, err := s.products.Get(id)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			l.Error().Err(err).Msg("failed to get products")
			return nil, fmt.Errorf("failed to get products: %w", ErrRequestFailed)
		}
		return nil, fmt.Errorf("product not found: %w", ErrNotFound)
	}

	return product, nil
}
