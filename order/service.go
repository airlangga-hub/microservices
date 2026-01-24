package main

import (
	"context"
	"errors"
	"log"

	accpb "github.com/airlangga-hub/microservices/order/account_pb"
	catpb "github.com/airlangga-hub/microservices/order/catalog_pb"
)

type Service interface {
	PostOrder(ctx context.Context, accountEmail string, accountID int32, products []OrderedProduct) (Order, error)
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error)
}

type service struct {
	repository    Repository
	accountClient accpb.AccountServiceClient
	catalogClient catpb.CatalogServiceClient
	publisher     *Publisher
}

func NewService(r Repository, accountClient accpb.AccountServiceClient, catalogClient catpb.CatalogServiceClient, publisher *Publisher) Service {
	return &service{
		repository:    r,
		accountClient: accountClient,
		catalogClient: catalogClient,
		publisher:     publisher,
	}
}

func (s *service) PostOrder(ctx context.Context, accountEmail string, accountID int32, products []OrderedProduct) (Order, error) {
	_, err := s.accountClient.GetAccount(ctx, &accpb.GetAccountRequest{AccountId: accountID})
	if err != nil {
		return Order{}, err
	}

	productIDs := make([]string, len(products))
	mapIdQty := make(map[string]int32)

	for i, p := range products {
		productIDs[i] = p.ID
		mapIdQty[p.ID] = p.Quantity
	}

	catalogRes, err := s.catalogClient.GetProducts(
		ctx,
		&catpb.GetProductsRequest{
			Offset: 0,
			Limit:  0,
			Ids:    productIDs,
			Query:  "",
		},
	)
	if err != nil {
		return Order{}, err
	}

	orderedProducts := make([]OrderedProduct, len(catalogRes.Products))
	var totalPrice int64

	for i, p := range catalogRes.Products {
		if qty, exist := mapIdQty[p.Id]; exist {
			orderedProducts[i] = OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    qty,
			}
			totalPrice += p.Price * int64(qty)
		}
	}

	if len(orderedProducts) != len(products) {
		log.Println("ERROR: order service PostOrder (check length): len(orderedProducts) != len(products)")
		return Order{}, errors.New("one or more products not found")
	}

	order := Order{
		AccountID:  accountID,
		Products:   orderedProducts,
		TotalPrice: totalPrice,
	}

	order, err = s.repository.CreateOrder(ctx, order)
	if err != nil {
		return Order{}, err
	}

	s.publisher.Publish(OrderMessage{
		ID:           order.ID,
		AccountEmail: accountEmail,
		Products:     order.Products,
		TotalPrice:   order.TotalPrice,
		CreatedAt:    order.CreatedAt,
	})

	return order, nil
}

func (s *service) GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error) {
	_, err := s.accountClient.GetAccount(ctx, &accpb.GetAccountRequest{AccountId: accountID})
	if err != nil {
		return nil, err
	}

	// the products inside each order only contain product_id and quantity
	// we need to get the name, description, price from catpb
	orders, err := s.repository.GetOrdersByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	productIdSet := make(map[string]struct{})

	for _, order := range orders {
		for _, product := range order.Products {
			productIdSet[product.ID] = struct{}{}
		}
	}

	productIDs := make([]string, 0, len(productIdSet))

	for productID := range productIdSet {
		productIDs = append(productIDs, productID)
	}

	catalogRes, err := s.catalogClient.GetProducts(
		ctx,
		&catpb.GetProductsRequest{
			Offset: 0,
			Limit:  0,
			Ids:    productIDs,
			Query:  "",
		},
	)
	if err != nil {
		return nil, err
	}

	mapCatalogProducts := make(map[string]*catpb.Product)

	for _, cp := range catalogRes.Products {
		mapCatalogProducts[cp.Id] = cp
	}

	for _, order := range orders {
		orderedProducts := make([]OrderedProduct, len(order.Products))
		for j, product := range order.Products {
			if cp, exist := mapCatalogProducts[product.ID]; exist {
				orderedProducts[j] = OrderedProduct{
					ID:          cp.Id,
					Name:        cp.Name,
					Description: cp.Description,
					Price:       cp.Price,
					Quantity:    product.Quantity,
				}
			}
		}

		if len(order.Products) != len(orderedProducts) {
			log.Println("ERROR: order service GetOrdersByAccountID (check length): len(order.Products) != len(orderedProducts)")
			return nil, errors.New("failed to find orders, one or more products not found")
		}

		order.Products = orderedProducts
	}

	return orders, nil
}
