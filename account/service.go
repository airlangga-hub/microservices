package main

import (
	"context"
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

type Service interface {
	SignUp(ctx context.Context, email, password string) (Account, error)
	Login(ctx context.Context, email, password string) (Account, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) SignUp(ctx context.Context, email, password string) (Account, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println("ERROR account service (CreateHash): ", err)
		return Account{}, errors.New("error creating account")
	}

	account, err := s.repository.CreateAccount(ctx, email, hashedPassword)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func (s *service) Login(ctx context.Context, email, password string) (Account, error) {
	account, err := s.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return Account{}, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, account.HashedPassword)
	if err != nil {
		log.Println("ERROR account service (ComparePasswordAndHash): ", err)
		return Account{}, errors.New("login error")
	}
	if !match {
		return Account{}, errors.New("incorrect password")
	}

	return account, nil
}
