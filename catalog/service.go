package main

import "context"

type Service interface {
	CreateProduct(ctx context.Context, name, description string, price int64) (Product, error)
	GetProductByID(ctx context.Context, id string) (Product, error)
	GetProducts(ctx context.Context, offset, limit int32) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, offset, limit int32) ([]Product, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateProduct(ctx context.Context, name, description string, price int64) (Product, error) {
	return s.repository.CreateProduct(ctx, productDocument{Name: name, Description: description, Price: price})
}

func (s *service) GetProductByID(ctx context.Context, id string) (Product, error) {
	return s.repository.GetProductByID(ctx, id)
}

func (s *service) GetProducts(ctx context.Context, offset, limit int32) ([]Product, error) {
	if limit > 100 || (offset == 0 && limit == 0) {
		limit = 100
	}

	return s.repository.ListProducts(ctx, offset, limit)
}

func (s *service) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.ListProductsWithIDs(ctx, ids)
}

func (s *service) SearchProducts(ctx context.Context, query string, offset, limit int32) ([]Product, error) {
	if limit > 100 || (offset == 0 && limit == 0) {
		limit = 100
	}

	return s.repository.SearchProducts(ctx, query, offset, limit)
}
