package account

import (
	"context"

	"github.com/airlangga-hub/microservices/services/account/pb"
)

type Server struct {
	svc Service
}

func (s *Server) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	account, err := s.svc.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	
	return &pb.PostAccountResponse{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *Server) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.svc.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	
	return &pb.GetAccountResponse{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *Server) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.svc.GetAccounts(ctx, r.Offset, r.Limit)
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