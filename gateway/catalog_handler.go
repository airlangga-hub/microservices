package main

import (
	"context"
	"net/http"
	"time"

	catpb "github.com/airlangga-hub/microservices/gateway/catalog_pb"
)

func (cfg *Config) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	res, err := cfg.CatalogClient.GetProducts(ctx, &catpb.GetProductsRequest{})
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
			Price:       product.Price,
		}
	}

	respondWithJSON(
		w,
		http.StatusOK,
		struct {
			Products []Product
		}{
			Products: products,
		},
	)
}
