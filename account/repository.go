package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close() error
	CreateAccount(ctx context.Context, email, hashedPassword string) (Account, error)
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	GetAccountByID(ctx context.Context, accountID int32) (Account, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(dbUrl string) (Repository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("ERROR: account repo NewRepository (sql.Open): ", err)
		return nil, errors.New("error connecting to db")
	}

	if err := db.Ping(); err != nil {
		log.Println("ERROR: account repo NewRepository (db.Ping): ", err)
		return nil, errors.New("error pinging db")
	}

	return &repository{db}, nil
}

func (r *repository) Close() error {
	if err := r.db.Close(); err != nil {
		log.Println("ERROR: account repo Close: ", err)
		return errors.New("error closing db")
	}
	return nil
}

func (r *repository) CreateAccount(ctx context.Context, email, hashedPassword string) (Account, error) {
	a := Account{}

	if err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO accounts
		(email, hashed_password)
		VALUES
		($1, $2)
		RETURNING id, email, type;`,
		email, hashedPassword,
	).Scan(&a.ID, &a.Email, &a.Type); err != nil {
		log.Println("ERROR: account repo CreateAccount: ", err)
		return Account{}, errors.New("error creating account")
	}

	return a, nil
}

func (r *repository) GetAccountByEmail(ctx context.Context, email string) (Account, error) {
	account := Account{}

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			email,
			hashed_password,
			type
		FROM accounts
		WHERE email = $1;`,
		email,
	).Scan(
		&account.ID,
		&account.Email,
		&account.HashedPassword,
		&account.Type,
	); err != nil {
		log.Println("ERROR: account repo GetAccountByEmail: ", err)
		return Account{}, errors.New("unauthorized user")
	}

	return account, nil
}

func (r *repository) GetAccountByID(ctx context.Context, accountID int32) (Account, error) {
	account := Account{}

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			email,
			type
		FROM accounts
		WHERE id = $1;`,
		accountID,
	).Scan(
		&account.ID,
		&account.Email,
		&account.Type,
	); err != nil {
		log.Println("ERROR: account repo GetAccountByID: ", err)
		return Account{}, errors.New("account not found")
	}

	return account, nil
}
