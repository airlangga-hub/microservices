package main

import (
	"context"
	"errors"
	"log"

	"github.com/airlangga-hub/microservices/order/pb"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	Svc Service
}

func (s *Server) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	products := make([]OrderedProduct, len(r.Products))
	for i, p := range r.Products {
		products[i] = OrderedProduct{
			ID:       p.Id,
			Quantity: p.Quantity,
		}
	}

	order, err := s.Svc.PostOrder(ctx, r.AccountEmail, r.AccountId, products)
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
			TotalPrice: order.TotalPrice,
			CreatedAt:  createdAt,
		},
	}, nil
}

func (s *Server) GetOrdersByAccountID(ctx context.Context, r *pb.GetOrdersByAccountIDRequest) (*pb.GetOrdersByAccountIDResponse, error) {
	orders, err := s.Svc.GetOrdersByAccountID(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}

	pbOrders := make([]*pb.Order, len(orders))

	for i, order := range orders {
		pbProducts := make([]*pb.OrderedProduct, len(order.Products))
		for j, product := range order.Products {
			pbProducts[j] = &pb.OrderedProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			}
		}
		
		createdAt, err := order.CreatedAt.MarshalBinary()
		if err != nil {
			log.Println("ERROR: order server GetOrdersByAccountID (MarshalBinary): ", err)
			return nil, errors.New("error finding orders")
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
