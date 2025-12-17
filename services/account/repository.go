package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	CreateAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id uint) (Account, error)
	ListAccounts(ctx context.Context, offset int, limit int) ([]*Account, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(dsn string) (Repository, error) {
	db, err := sql.Open()
}

func (r *repository) CreateAccount(ctx context.Context, a Account) error {
	
}

func (r *repository) GetAccountByID(ctx context.Context, id uint) (Account, error) {
	
}

func (r *repository) ListAccounts(ctx context.Context, offset int, limit int) ([]*Account, error) {
	
}