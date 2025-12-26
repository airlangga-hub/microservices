package order

import (
	"context"
	"fmt"
	"net"

	"github.com/airlangga-hub/microservices/services/account"
	"github.com/airlangga-hub/microservices/services/catalog"
	"github.com/airlangga-hub/microservices/services/order/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	Svc           Service
	AccountClient *account.Client
	CatalogClient *catalog.Client
}

func ListenGrpc(service Service, port string) error {
	accountClient, err := account.NewClient()
	if err != nil {
		return fmt.Errorf("ERROR: order server ListenGrpc (account.NewClient): %v", err)
	}
	defer accountClient.Conn.Close()

	catalogClient, err := catalog.NewClient()
	if err != nil {
		return fmt.Errorf("ERROR: order server ListenGrpc (catalog.NewClient): %v", err)
	}
	defer catalogClient.Conn.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("ERROR: order server ListenGrpc (net.Listen): %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterOrderServiceServer(
		s,
		&Server{
			Svc:           service,
			AccountClient: accountClient,
			CatalogClient: catalogClient,
		},
	)

	return s.Serve(lis)
}

func (s *Server) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.AccountClient.GetAccount(ctx, r.AccountId)
}

func (s *Server) GetOrdersByAccountID(ctx context.Context, r *pb.GetOrdersByAccountIDRequest) (*pb.GetOrdersByAccountIDResponse, error)
