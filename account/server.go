package main

import (
	"context"

	"github.com/airlangga-hub/microservices/account/pb"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
	Svc Service
}

func (s *Server) SignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.PostAccountResponse, error) {
	account, err := s.Svc.SignUp(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Account{Email: account.Email, }, nil
}

func (s *Server) Login(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.Svc.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}
