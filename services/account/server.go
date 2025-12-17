package account

import "context"

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