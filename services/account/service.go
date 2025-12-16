package account

import (
	"context"

	"github.com/airlangga-hub/microservices/services/account/domain"
)

type Service interface {
	PostAccount(ctx context.Context, name string) error
	GetAccount(ctx context.Context, id uint) (domain.Account, error)
	GetAccounts(ctx context.Context, offset, limit int) ([]*domain.Account, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) PostAccount(ctx context.Context, name string) error {
	account := domain.Account{Name: name}
	return s.repository.CreateAccount(ctx, account)
}

func (s *service) GetAccount(ctx context.Context, id uint) (domain.Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *service) GetAccounts(ctx context.Context, offset, limit int) ([]*domain.Account, error) {
	return s.repository.ListAccounts(ctx, offset, limit)
}