package account

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close() error
	CreateAccount(ctx context.Context, a Account) (Account, error)
	GetAccountByID(ctx context.Context, id int32) (Account, error)
	ListAccounts(ctx context.Context, offset, limit int32) ([]Account, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(dbUrl string) (Repository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("ERROR: account repo NewRepository: ", err)
		return nil, errors.New("error connecting to db")
	}

	if err := db.Ping(); err != nil {
		log.Println("ERROR: account repo NewRepository: ", err)
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

func (r *repository) CreateAccount(ctx context.Context, a Account) (Account, error) {
	if err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO accounts (name)
		VALUES ($1)
		RETURNING id, name;`,
		a.Name,
	).Scan(&a.ID, &a.Name); err != nil {
		log.Println("ERROR: account repo CreateAccount: ", err)
		return Account{}, errors.New("error creating account")
	}

	return a, nil
}

func (r *repository) GetAccountByID(ctx context.Context, id int32) (Account, error) {
	account := Account{}

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			name
		FROM accounts
		WHERE id = $1;`,
		id,
	).Scan(&account.ID, &account.Name); err != nil {
		log.Println("ERROR: account repo GetAccountByID: ", err)
		return Account{}, errors.New("error getting account by id")
	}

	return account, nil
}

func (r *repository) ListAccounts(ctx context.Context, offset, limit int32) ([]Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			id,
			name
		FROM accounts
		ORDER BY id DESC
		OFFSET $1
		LIMIT $2;`,
		offset,
		limit)
	if err != nil {
		log.Println("ERROR: account repo ListAccounts: ", err)
		return nil, errors.New("error listing accounts")
	}

	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		a := Account{}
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			log.Println("ERROR: account repo ListAccounts: ", err)
			return nil, errors.New("error scanning current row")
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		log.Println("ERROR: account repo ListAccounts: ", err)
		return nil, errors.New("error iterating rows")
	}

	return accounts, nil
}
