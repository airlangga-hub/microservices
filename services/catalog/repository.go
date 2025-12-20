package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

type Repository interface {
	CreateProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (Product, error)
	ListProducts(ctx context.Context, offset, limit int32) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, offset, limit int32) ([]Product, error)
}

type repository struct {
	client *elasticsearch.Client
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

const ESIndex = "catalog"

type ESresponse struct {
	Hits struct {
		Hits []struct {
			Source Product `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewRepository() (Repository, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Println("ERROR: catalog repo NewRepository: ", err)
		return nil, errors.New("error creating elastic search client")
	}

	// create index
	client.Indices.Create(ESIndex)

	return &repository{client}, nil
}

func (r *repository) CreateProduct(ctx context.Context, p Product) error {
	product, err := json.Marshal(p)
	if err != nil {
		log.Println("ERROR: catalog repo CreateProduct: ", err)
		return errors.New("error marshaling product")
	}

	req := esapi.IndexRequest{
		Index:   ESIndex,
		Body:    bytes.NewReader(product),
		Refresh: "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo CreateProduct: ", err)
		return errors.New("error creating product in elastic search")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo CreateProduct status=%d body=%s", res.StatusCode, body)
		return errors.New("error creating product in elastic search")
	}

	var response struct {
		ID string `json:"_id"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo CreateProduct: decode ID error", err)
		return errors.New("error decoding generated ID")
	}

	log.Println("CREATED: created product with id: ", response.ID)

	return nil
}

func (r *repository) GetProductByID(ctx context.Context, id string) (Product, error) {
	req := esapi.GetRequest{
		Index:      ESIndex,
		DocumentID: id,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo GetProductByID: ", err)
		return Product{}, errors.New("error getting product by id in elastic search")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo GetProductByID status=%d body=%s", res.StatusCode, body)
		return Product{}, errors.New("error getting product by id in elastic search")
	}

	var response struct {
		Source Product `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo GetProductByID: ", err)
		return Product{}, errors.New("error decoding get product by id response")
	}

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
		log.Println("ERROR: catalog repo ListProducts: ", err)
		return nil, errors.New("error marshaling query for ListProducts")
	}

	req := esapi.SearchRequest{
		Index: []string{ESIndex},
		Body:  bytes.NewReader(esQuery),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo ListProducts: ", err)
		return nil, errors.New("error listing products")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo ListProducts: status=%d, body=%s", res.StatusCode, body)
		return nil, errors.New("error listing products")
	}

	var response ESresponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo ListProducts: ", err)
		return nil, errors.New("error decoding ListProducts response")
	}

	products := []Product{}

	for _, hit := range response.Hits.Hits {
		products = append(products, hit.Source)
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
		log.Println("ERROR: catalog repo ListProductsWithIDs: ", err)
		return nil, errors.New("error marshaling query for ListProductsWithIDs")
	}

	req := esapi.SearchRequest{
		Index: []string{ESIndex},
		Body:  bytes.NewReader(esQuery),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo ListProductsWithIDs: ", err)
		return nil, errors.New("error searching products by IDs")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo ListProductsWithIDs: status=%d, body=%s", res.StatusCode, body)
		return nil, errors.New("error listing products by IDs")
	}

	var response ESresponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo ListProductsWithIDs: ", err)
		return nil, errors.New("error decoding products by IDs response")
	}

	products := []Product{}

	for _, hit := range response.Hits.Hits {
		products = append(products, hit.Source)
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
		log.Println("ERROR: catalog repo SearchProducts: ", err)
		return nil, errors.New("error marshaling search products query")
	}

	req := esapi.SearchRequest{
		Index: []string{ESIndex},
		Body:  bytes.NewReader(esQuery),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("ERROR: catalog repo SearchProducts: ", err)
		return nil, errors.New("error searching products")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		log.Printf("ERROR: catalog repo SearchProducts: status=%d, body=%s", res.StatusCode, body)
		return nil, errors.New("elasticsearch error searching products")
	}

	var response ESresponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("ERROR: catalog repo SearchProducts: decode error", err)
		return nil, errors.New("error decoding search results")
	}

	products := []Product{}

	for _, hit := range response.Hits.Hits {
		products = append(products, hit.Source)
	}

	return products, nil
}
