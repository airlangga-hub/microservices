package account

import (
	"context"
)

type Account struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	PostAccount(ctx context.Context, name string) error
	GetAccount(ctx context.Context, id int32) (Account, error)
	GetAccounts(ctx context.Context, offset, limit int32) ([]Account, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) PostAccount(ctx context.Context, name string) error {
	return s.repository.CreateAccount(ctx, Account{Name: name})
}

func (s *service) GetAccount(ctx context.Context, id int32) (Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *service) GetAccounts(ctx context.Context, offset, limit int32) ([]Account, error) {
	if limit > 100 || (offset == 0 && limit == 0) {
		limit = 100
	}
	
	return s.repository.ListAccounts(ctx, offset, limit)
}
