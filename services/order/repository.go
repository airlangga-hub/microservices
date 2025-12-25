package order

import "context"

type Order struct {
	ID        int32 `json:"id"`
	AccountID int32 `json:"id"`
	TotalPrice 
}

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, o Order) (Order, error)
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]Order, error)
}
