package order

import "context"

type Service interface {
	PostOrder(ctx context.Context, accountID int32, products []OrderedProduct) (Order, error)
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error)
}

type service struct {
	repository Repository
}

func (s *service) PostOrder(ctx context.Context, accountID int32, products []OrderedProduct) (Order, error)

func (s *service) GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error)
