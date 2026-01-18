package main

import (
	"context"
	"errors"
	"log"

	accpb "github.com/airlangga-hub/microservices/order/account_pb"
	catpb "github.com/airlangga-hub/microservices/order/catalog_pb"
	"github.com/airlangga-hub/microservices/order/pb"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	Svc           Service
	AccountClient accpb.AccountServiceClient
	CatalogClient catpb.CatalogServiceClient
}


func (s *Server) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.AccountClient.GetAccount(ctx, &accpb.GetAccountRequest{AccountId: r.AccountId})
	if err != nil {
		return nil, err
	}

	productIDs := make([]string, len(r.Products))
	mapIdQty := make(map[string]int32)

	for i, p := range r.Products {
		productIDs[i] = p.Id
		mapIdQty[p.Id] = p.Quantity
	}

	catalogRes, err := s.CatalogClient.GetProducts(
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

	orderedProducts := make([]OrderedProduct, len(catalogRes.Products))

	for i, p := range catalogRes.Products {
		if qty, exist := mapIdQty[p.Id]; exist {
			orderedProducts[i] = OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    qty,
			}
		}
	}

	if len(orderedProducts) != len(r.Products) {
		log.Println("ERROR: order server PostOrder (check length): ", err)
		return nil, errors.New("one or more products not found")
	}

	order, err := s.Svc.PostOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		return nil, err
	}

	pbProducts := make([]*pb.OrderedProduct, len(order.Products))

	for i, p := range order.Products {
		pbProducts[i] = &pb.OrderedProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		}

	}

	createdAt, err := order.CreatedAt.MarshalBinary()
	if err != nil {
		log.Println("ERROR: order server PostOrder (MarshalBinary): ", err)
		return nil, errors.New("error creating order")
	}

	return &pb.PostOrderResponse{
		Order: &pb.Order{
			Id:         order.ID,
			AccountId:  order.AccountID,
			Products:   pbProducts,
			TotalPrice: order.TotalPrice / 100,
			CreatedAt:  createdAt,
		},
	}, nil
}

func (s *Server) GetOrdersByAccountID(ctx context.Context, r *pb.GetOrdersByAccountIDRequest) (*pb.GetOrdersByAccountIDResponse, error) {
	_, err := s.AccountClient.GetAccount(ctx, &accpb.GetAccountRequest{AccountId: r.AccountId})
	if err != nil {
		return nil, err
	}

	// the products inside each order only contains product_id and quantity
	// we need to get the name, description, price from catpb
	orders, err := s.Svc.GetOrdersByAccountID(ctx, r.AccountId)
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

	catalogRes, err := s.CatalogClient.GetProducts(
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

	pbOrders := make([]*pb.Order, len(orders))

	for i, order := range orders {
		pbProducts := make([]*pb.OrderedProduct, len(order.Products))
		for j, product := range order.Products {
			if cp, exist := mapCatalogProducts[product.ID]; exist {
				pbProducts[j] = &pb.OrderedProduct{
					Id:          cp.Id,
					Name:        cp.Name,
					Description: cp.Description,
					Price:       cp.Price,
					Quantity:    product.Quantity,
				}
			}
		}

		if len(order.Products) != len(pbProducts) {
			log.Println("ERROR: order server GetOrdersByAccountID (check length): ", err)
			return nil, errors.New("failed to find orders")
		}

		createdAt, err := order.CreatedAt.MarshalBinary()
		if err != nil {
			log.Println("ERROR: order server GetOrdersByAccountID (MarshalBinary): ", err)
			return nil, errors.New("failed to find orders")
		}

		pbOrders[i] = &pb.Order{
			Id:         order.ID,
			AccountId:  order.AccountID,
			Products:   pbProducts,
			TotalPrice: order.TotalPrice,
			CreatedAt:  createdAt,
		}
	}

	return &pb.GetOrdersByAccountIDResponse{
		Orders: pbOrders,
	}, nil
}
