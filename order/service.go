package main

import "context"

type Service interface {
	PostOrder(ctx context.Context, accountID int32, products []OrderedProduct) (Order, error)
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) PostOrder(ctx context.Context, accountID int32, products []OrderedProduct) (Order, error) {
	order := Order{
		AccountID: accountID,
		Products:  products,
	}

	for _, p := range products {
		order.TotalPrice += p.Price * int64(p.Quantity)
	}

	return s.repository.CreateOrder(ctx, order)
}

func (s *service) GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error) {
	return s.repository.GetOrdersByAccountID(ctx, accountID)
}
