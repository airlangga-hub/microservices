package order

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
)

type OrderedProduct struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Quantity    int32  `json:"quantity"`
}

type Order struct {
	ID         int32            `json:"id"`
	AccountID  int32            `json:"account_id"`
	Products   []OrderedProduct `json:"products"`
	TotalPrice int64            `json:"total_price"`
	CreatedAt  time.Time        `json:"created_at"`
}

type Repository interface {
	Close() error
	CreateOrder(ctx context.Context, o Order) (Order, error)
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]Order, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(dbUrl string) (Repository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("ERROR: order repo NewRepository: ", err)
		return nil, errors.New("error opening postgres")
	}

	if err := db.Ping(); err != nil {
		log.Println("ERROR: order repo NewRepository: ", err)
		return nil, errors.New("error pinging db")
	}

	return &repository{db}, nil
}

func (r *repository) Close() error {
	if err := r.db.Close(); err != nil {
		log.Println("ERROR: order repo Close: ", err)
		return errors.New("error closing db")
	}
	return nil
}

func (r *repository) CreateOrder(ctx context.Context, o Order) (Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		log.Println("ERROR: order repo CreateOrder (tx init): ", err)
		return Order{}, errors.New("error creating order")
	}

	// insert order
	if err = tx.QueryRowContext(
		ctx,
		`INSERT INTO orders (account_id, total_price)
		VALUES ($1, $2)
		RETURNING id, created_at;`,
		o.AccountID, o.TotalPrice,
	).Scan(&o.ID, &o.CreatedAt); err != nil {
		log.Println("ERROR: order repo CreateOrder (insert order): ", err)
		return Order{}, errors.New("error creating order")
	}

	// insert order products
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	defer stmt.Close()
	if err != nil {
		log.Println("ERROR: order repo CreateOrder (stmt prepare): ", err)
		return Order{}, errors.New("error creating order")
	}

	for _, p := range o.Products {
		_, err := stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			log.Println("ERROR: order repo CreateOrder (insert order products): ", err)
			return Order{}, errors.New("error creating order")
		}
	}

	// flush
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Println("ERROR: order repo CreateOrder (flush): ", err)
		return Order{}, errors.New("error creating order")
	}

	if err = tx.Commit(); err != nil {
		log.Println("ERROR: order repo CreateOrder (tx commit): ", err)
		return Order{}, errors.New("error creating order")
	}

	return o, nil
}

func (r *repository) GetOrdersByAccountID(ctx context.Context, accountID int32) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			id,
			account_id,
			total_price,
			created_at
		FROM orders
		WHERE account_id = $1;`,
		accountID,
	)

	if err != nil {
		log.Println("ERROR: order repo GetOrdersByAccountID: ", err)
		return nil, errors.New("error finding account's orders")
	}

	defer rows.Close()

	orders := []Order{}

	for rows.Next() {
		o := Order{}
		if err := rows.Scan(&o.ID, &o.AccountID, &o.TotalPrice, &o.CreatedAt); err != nil {
			log.Println("ERROR: order repo GetOrdersByAccountID: ", err)
			return nil, errors.New("error finding account's orders")
		}
		orders = append(orders, o)
	}

	return orders, nil
}
