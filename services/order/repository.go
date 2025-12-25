package order

import (
	"context"
	"time"
)

type Order struct {
	ID         int32     `json:"id"`
	AccountID  int32     `json:"account_id"`
	TotalPrice int64     `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, o Order) (Order, error)
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]Order, error)
}
