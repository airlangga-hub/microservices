package account

import (
	"context"
	"errors"
	"log"

	"github.com/airlangga-hub/microservices/services/account/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository interface {
	CreateAccount(ctx context.Context, a domain.Account) error
	GetAccountByID(ctx context.Context, id uint) (domain.Account, error)
	ListAccounts(ctx context.Context, skip uint, take uint) ([]*domain.Account, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Println("ERROR: accounts repo --> NewRepository: ", err)
		return nil, errors.New("database error")
	}

	return &repository{db}, nil
}

func (r *repository) CreateAccount(ctx context.Context, a domain.Account) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		log.Println("ERROR: accounts repo --> CreateAccount: ", err)
		return errors.New("error creating account")
	}
	
	return nil
}

func (r *repository) GetAccountByID(ctx context.Context, id uint) (domain.Account, error) {
	account := domain.Account{}
	
	if err := r.db.WithContext(ctx).First(&account, "id=?", id).Error; err != nil {
		log.Println("ERROR: accounts repo --> GetAccountByID: ", err)
		return domain.Account{}, errors.New("account not found")
	}
	
	return account, nil
}

func (r *repository) ListAccounts(ctx context.Context, skip uint, take uint) ([]*domain.Account, error) {
	accounts := []*domain.Account{}
	
	if err := r.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		log.Println("ERROR: accounts repo --> ListAccounts: ", err)
		return nil, errors.New("no accounts found")
	}
	
	if len(accounts) == 0 {
		return nil, errors.New("no accounts found")
	}
	
	return accounts, nil
}