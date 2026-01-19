package main

import (
	"context"
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

type Service interface {
	SignUp(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetAccount(ctx context.Context, accountID int32) (Account, error)
}

type service struct {
	repository Repository
	secret     []byte
}

func NewService(r Repository, secret []byte) Service {
	return &service{repository: r, secret: secret}
}

func (s *service) SignUp(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println("ERROR account service (CreateHash): ", err)
		return "", errors.New("error creating account")
	}

	account, err := s.repository.CreateAccount(ctx, email, hashedPassword)
	if err != nil {
		return "", err
	}

	return MakeJWT(account.ID, account.Type, account.Email, s.secret)
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	account, err := s.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	match, err := argon2id.ComparePasswordAndHash(password, account.HashedPassword)
	if err != nil {
		log.Println("ERROR account service (ComparePasswordAndHash): ", err)
		return "", errors.New("login error")
	}
	if !match {
		return "", errors.New("incorrect password")
	}

	return MakeJWT(account.ID, account.Type, account.Email, s.secret)
}

func (s *service) GetAccount(ctx context.Context, accountID int32) (Account, error) {
	return s.repository.GetAccountByID(ctx, accountID)
}
