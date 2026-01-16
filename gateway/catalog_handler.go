package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	catpb "github.com/airlangga-hub/microservices/gateway/catalog_pb"
)

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	res, err := h.CatalogClient.GetProducts(ctx, &catpb.GetProductsRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(res.Products))

	for i, product := range res.Products {
		products[i] = Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price / 100,
		}
	}

	respondWithJSON(
		w,
		http.StatusOK,
		struct {
			Products []Product `json:"products"`
		}{
			Products: products,
		},
	)
}

func (h *Handler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	res, err := h.CatalogClient.GetProducts(ctx, &catpb.GetProductsRequest{Query: request.Query})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(res.Products))

	for i, product := range res.Products {
		products[i] = Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price / 100,
		}
	}

	respondWithJSON(
		w,
		http.StatusOK,
		struct {
			Products []Product `json:"products"`
		}{
			Products: products,
		},
	)
}
