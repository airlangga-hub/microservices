package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	orderpb "github.com/airlangga-hub/microservices/gateway/order_pb"
)

func (cfg *Config) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	account, ok := r.Context().Value(accountKey).(Account)
	if !ok {
		http.Error(w, "unauthorized user", http.StatusUnauthorized)
		return
	}

	var request struct {
		Products []Product `json:"products"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pbProducts := make([]*orderpb.OrderedProduct, len(request.Products))

	for i, p := range request.Products {
		pbProducts[i] = &orderpb.OrderedProduct{
			Id:       p.ID,
			Quantity: p.Quantity,
		}
	}

	order, err := cfg.OrderClient.PostOrder(ctx, &orderpb.PostOrderRequest{AccountId: account.ID, Products: pbProducts})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var t time.Time
	if err := t.UnmarshalBinary(order.Order.CreatedAt); err != nil {
		log.Println("ERROR CreateOrder (UnmarshalBinary): ", err)
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(order.Order.Products))

	for i, p := range order.Order.Products {
		products[i] = Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		}
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		struct {
			OrderID    int32     `json:"order_id"`
			AccountID  int32     `json:"account_id"`
			Products   []Product `json:"products"`
			TotalPrice int64     `json:"total_price"`
			CreatedAt  time.Time `json:"created_at"`
		}{
			OrderID:    order.Order.Id,
			AccountID:  order.Order.AccountId,
			Products:   products,
			TotalPrice: order.Order.TotalPrice,
			CreatedAt:  t,
		},
	)
}
