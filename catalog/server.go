package main

import (
	"context"

	"github.com/airlangga-hub/microservices/catalog/pb"
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
	Svc Service
}

func (s *Server) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	product, err := s.Svc.CreateProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}

func (s *Server) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.Svc.GetProductByID(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price / 100,
		},
	}, nil
}

func (s *Server) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var products []Product
	var err error

	if r.Query != "" {
		products, err = s.Svc.SearchProducts(ctx, r.Query, r.Offset, r.Limit)
	} else if len(r.Ids) > 0 {
		products, err = s.Svc.GetProductsByIDs(ctx, r.Ids)
	} else {
		products, err = s.Svc.GetProducts(ctx, r.Offset, r.Limit)
	}

	if err != nil {
		return nil, err
	}

	pbProducts := make([]*pb.Product, len(products))

	for i, p := range products {
		pbProducts[i] = &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price / 100,
		}
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}
