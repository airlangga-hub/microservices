package main

import (
	"context"

	"github.com/airlangga-hub/microservices/account/pb"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
	Svc Service
}

func (s *Server) SignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.Token, error) {
	token, err := s.Svc.SignUp(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Token{Jwt: token}, nil
}

func (s *Server) Login(ctx context.Context, r *pb.LoginRequest) (*pb.Token, error) {
	token, err := s.Svc.Login(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Token{Jwt: token}, nil
}

func (s *Server) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.Svc.GetAccount(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		AccountId: account.ID,
		Email:     account.Email,
	}, nil
}
