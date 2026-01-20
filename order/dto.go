package main

import "time"

type OrderedProduct struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Quantity    int32  `json:"quantity"`
}

type Order struct {
	ID         int32            `json:"id"`
	AccountID  int32            `json:"account_id"`
	Products   []OrderedProduct `json:"products"`
	TotalPrice int64            `json:"total_price"`
	CreatedAt  time.Time        `json:"created_at"`
}
