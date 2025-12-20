package account

import (
	"context"
	"log"
	"net"

	"github.com/airlangga-hub/microservices/services/account/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
	Svc Service
}

func ListenGrpc(service Service, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR: account server ListenGrpc: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterAccountServiceServer(s, &Server{Svc: service})

	return s.Serve(lis)
}

func (s *Server) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	account, err := s.Svc.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *Server) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.Svc.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *Server) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.Svc.GetAccounts(ctx, r.Offset, r.Limit)
	if err != nil {
		return nil, err
	}

	a := []*pb.Account{}

	for _, account := range accounts {
		pbAccount := &pb.Account{Id: account.ID, Name: account.Name}
		a = append(a, pbAccount)
	}

	return &pb.GetAccountsResponse{Accounts: a}, nil
}
