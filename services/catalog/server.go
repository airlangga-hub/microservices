package catalog

import (
	"context"
	"log"
	"net"

	"github.com/airlangga-hub/microservices/services/catalog/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
	Svc Service
}

func ListenGrpc(service Service, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR: account server ListenGrpc: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCatalogServiceServer(s, &Server{Svc: service})

	return s.Serve(lis)
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
			Price:       product.Price,
		},
	}, nil
}

func (s *Server) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	products := []Product{}
	var err error

	if r.Query != "" {
		products, err = s.Svc.SearchProducts(ctx, r.Query, r.Offset, r.Limit)
	} else if r.Ids != nil {
		products, err = s.Svc.GetProductsByIDs(ctx, r.Ids)
	} else {
		products, err = s.Svc.GetProducts(ctx, r.Offset, r.Limit)
	}

	if err != nil {
		return nil, err
	}

	pbProducts := []*pb.Product{}

	for _, p := range products {
		pbProducts = append(
			pbProducts,
			&pb.Product{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		)
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}
