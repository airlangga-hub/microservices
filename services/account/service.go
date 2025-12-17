package account

import (
	"context"
)

type Account struct {
	ID   int32   `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	PostAccount(ctx context.Context, name string) error
	GetAccount(ctx context.Context, id uint) (Account, error)
	GetAccounts(ctx context.Context, offset, limit int) ([]*Account, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) PostAccount(ctx context.Context, name string) error {
	account := Account{Name: name}
	return s.repository.CreateAccount(ctx, account)
}

func (s *service) GetAccount(ctx context.Context, id uint) (Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *service) GetAccounts(ctx context.Context, offset, limit int) ([]*Account, error) {
	return s.repository.ListAccounts(ctx, offset, limit)
}
