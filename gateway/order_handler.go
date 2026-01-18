package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	orderpb "github.com/airlangga-hub/microservices/gateway/order_pb"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	res, err := h.OrderClient.PostOrder(ctx, &orderpb.PostOrderRequest{AccountId: account.ID, Products: pbProducts})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var t time.Time
	if err := t.UnmarshalBinary(res.Order.CreatedAt); err != nil {
		log.Println("ERROR gateway CreateOrder (UnmarshalBinary): ", err)
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(res.Order.Products))

	for i, p := range res.Order.Products {
		products[i] = Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       float64(p.Price / 100),
			Quantity:    p.Quantity,
		}
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		Order{
			ID:         res.Order.Id,
			AccountID:  res.Order.AccountId,
			Products:   products,
			TotalPrice: float64(res.Order.TotalPrice / 100),
			CreatedAt:  t,
		},
	)
}

func (h *Handler) GetOrdersByAccountID(w http.ResponseWriter, r *http.Request) {
	account, ok := r.Context().Value(accountKey).(Account)
	if !ok {
		http.Error(w, "unauthorized user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	res, err := h.OrderClient.GetOrdersByAccountID(ctx, &orderpb.GetOrdersByAccountIDRequest{AccountId: account.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orders := make([]*Order, len(res.Orders))

	for i, o := range res.Orders {
		products := make([]Product, len(o.Products))
		for j, p := range o.Products {
			products[j] = Product{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       float64(p.Price / 100),
				Quantity:    p.Quantity,
			}
		}

		var t time.Time
		if err := t.UnmarshalBinary(o.CreatedAt); err != nil {
			log.Println("ERROR gateway GetOrdersByAccountID (UnmarshalBinary): ", err)
			http.Error(w, "some error happened", http.StatusInternalServerError)
			return
		}

		orders[i] = &Order{
			ID:         o.Id,
			AccountID:  o.AccountId,
			Products:   products,
			TotalPrice: float64(o.TotalPrice / 100),
			CreatedAt:  t,
		}
	}

	respondWithJSON(
		w,
		http.StatusOK,
		struct {
			Orders []*Order `json:"orders"`
		}{
			Orders: orders,
		},
	)
}
