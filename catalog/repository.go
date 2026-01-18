package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

type Repository interface {
	Close(ctx context.Context) error
	CreateProduct(ctx context.Context, p productDocument) (Product, error)
	GetProductByID(ctx context.Context, id string) (Product, error)
	ListProducts(ctx context.Context, offset, limit int32) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, offset, limit int32) ([]Product, error)
}

type repository struct {
	client *elasticsearch.Client
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	SellerID    int32  `json:"seller_id"`
}

const ESIndex = "catalog"

type ESresponse struct {
	Hits struct {
		Hits []struct {
			ID     string  `json:"_id"`
			Source Product `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewRepository() (Repository, error) {
	esUrl := os.Getenv("ELASTICSEARCH_URL")
	cfg := elasticsearch.Config{Addresses: []string{esUrl}}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// 1. Check if index exists
	_, err = esapi.IndicesExistsRequest{
		Index: []string{ESIndex},
	}.Do(context.Background(), client)

	if err != nil {
		client.Indices.Create(ESIndex)
	}

	return &repository{client}, nil
}

func (r *repository) Close(ctx context.Context) error {
	if err := r.client.Close(ctx); err != nil {
		log.Println("ERROR: catalog repo Close: ", err)
		return errors.New("error closing elastic search client")
	}

	return nil
}

func (r *repository) CreateProduct(ctx context.Context, p productDocument) (Product, error) {
	productDoc, err := json.Marshal(p)
	if err != nil {
		log.Println("ERROR: catalog repo CreateProduct (json.Marshal): ", err)
		return Product{}, errors.New("failed to create product")
	}

	req := esapi.IndexRequest{
		Index:   ESIndex,
		Body:    bytes.NewReader(productDoc),
		Refresh: "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo CreateProduct (req.Do): ", err)
		return Product{}, errors.New("failed to create product")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo CreateProduct status=%d body=%s\n", res.StatusCode, body)
		return Product{}, errors.New("failed to create product")
	}

	var response struct {
		ID     string  `json:"_id"`
		Source Product `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo CreateProduc (es response decode): ", err)
		return Product{}, errors.New("failed to create product")
	}

	response.Source.ID = response.ID
	response.Source.Name = p.Name
	response.Source.Description = p.Description
	response.Source.Price = p.Price

	return response.Source, nil
}

func (r *repository) GetProductByID(ctx context.Context, id string) (Product, error) {
	req := esapi.GetRequest{
		Index:      ESIndex,
		DocumentID: id,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo GetProductByID (req.Do): ", err)
		return Product{}, errors.New("failed to get product by id")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo GetProductByID status=%d body=%s\n", res.StatusCode, body)
		return Product{}, errors.New("failed to get product by id")
	}

	var response struct {
		Source Product `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo GetProductByID (es response decode): ", err)
		return Product{}, errors.New("failed to get product by id")
	}

	response.Source.ID = id

	return response.Source, nil
}

func (r *repository) ListProducts(ctx context.Context, offset, limit int32) ([]Product, error) {
	query := map[string]any{
		"from": offset,
		"size": limit,
		"query": map[string]any{
			"match_all": map[string]any{},
		},
	}

	esQuery, err := json.Marshal(query)
	if err != nil {
		log.Println("ERROR: catalog repo ListProducts (json.Marshal): ", err)
		return nil, errors.New("failed to find products")
	}

	req := esapi.SearchRequest{
		Index: []string{ESIndex},
		Body:  bytes.NewReader(esQuery),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo ListProducts (req.Do): ", err)
		return nil, errors.New("failed to find products")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo ListProducts: status=%d, body=%s\n", res.StatusCode, body)
		return nil, errors.New("failed to find products")
	}

	var response ESresponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo ListProducts (es response decode): ", err)
		return nil, errors.New("failed to find products")
	}

	products := make([]Product, len(response.Hits.Hits))

	for i, hit := range response.Hits.Hits {
		hit.Source.ID = hit.ID
		products[i] = hit.Source
	}

	return products, nil
}

func (r *repository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	if len(ids) == 0 {
		return []Product{}, nil
	}

	query := map[string]any{
		"query": map[string]any{
			"ids": map[string]any{
				"values": ids,
			},
		},
	}

	esQuery, err := json.Marshal(query)
	if err != nil {
		log.Println("ERROR: catalog repo ListProductsWithIDs (json.Marshal): ", err)
		return nil, errors.New("failed to find products by IDs")
	}

	req := esapi.SearchRequest{
		Index: []string{ESIndex},
		Body:  bytes.NewReader(esQuery),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo ListProductsWithIDs (req.Do): ", err)
		return nil, errors.New("failed to find products by IDs")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo ListProductsWithIDs: status=%d, body=%s\n", res.StatusCode, body)
		return nil, errors.New("failed to find products by IDs")
	}

	var response ESresponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo ListProductsWithIDs (es response decode): ", err)
		return nil, errors.New("failed to find products by IDs")
	}

	products := make([]Product, len(response.Hits.Hits))

	for i, hit := range response.Hits.Hits {
		hit.Source.ID = hit.ID
		products[i] = hit.Source
	}

	return products, nil
}

func (r *repository) SearchProducts(ctx context.Context, query string, offset, limit int32) ([]Product, error) {
	q := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  query,
				"fields": []string{"name^2", "description"},
			},
		},
		"from": offset,
		"size": limit,
	}

	esQuery, err := json.Marshal(q)
	if err != nil {
		log.Println("ERROR: catalog repo SearchProducts (json.Marshal): ", err)
		return nil, errors.New("failed to search products")
	}

	req := esapi.SearchRequest{
		Index: []string{ESIndex},
		Body:  bytes.NewReader(esQuery),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo SearchProducts (req.Do): ", err)
		return nil, errors.New("failed to search products")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo SearchProducts: status=%d, body=%s\n", res.StatusCode, body)
		return nil, errors.New("failed to search products")
	}

	var response ESresponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo SearchProducts (es response decode): ", err)
		return nil, errors.New("failed to search products")
	}

	products := make([]Product, len(response.Hits.Hits))

	for i, hit := range response.Hits.Hits {
		hit.Source.ID = hit.ID
		products[i] = hit.Source
	}

	return products, nil
}
