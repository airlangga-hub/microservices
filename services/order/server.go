package order

import (
	"context"
	"log"
	"net"

	"github.com/airlangga-hub/microservices/services/order/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	Svc Service
}

func ListenGrpc(service Service, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR: account server ListenGrpc: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterOrderServiceServer(s, &Server{Svc: service})

	return s.Serve(lis)
}

func (s *Server) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error)

func (s *Server) GetOrdersByAccountID(ctx context.Context, r *pb.GetOrdersByAccountIDRequest) (*pb.GetOrdersByAccountIDResponse, error)
