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
	GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(dbUrl string) (Repository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("ERROR: order repo NewRepository (sql.Open): ", err)
		return nil, errors.New("error opening postgres")
	}

	if err := db.Ping(); err != nil {
		log.Println("ERROR: order repo NewRepository (db.Ping): ", err)
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
		RETURNING
			id,
			created_at;`,
		o.AccountID, o.TotalPrice,
	).Scan(
		&o.ID,
		&o.CreatedAt,
	); err != nil {
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

func (r *repository) GetOrdersByAccountID(ctx context.Context, accountID int32) ([]*Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			o.id,
			o.account_id,
			o.total_price,
			o.created_at,
			op.product_id,
			op.quantity
		FROM orders o
		JOIN order_products op
		ON o.id = op.order_id
		WHERE account_id = $1
		ORDER BY o.id;`,
		accountID,
	)

	if err != nil {
		log.Println("ERROR: order repo GetOrdersByAccountID (r.db.QueryContext): ", err)
		return nil, errors.New("error finding account's orders")
	}

	defer rows.Close()

	ordersMap := map[int32]*Order{}

	for rows.Next() {
		var (
			id          int32
			account_id  int32
			total_price int64
			created_at  time.Time
			product_id  string
			quantity    int32
		)

		if err := rows.Scan(
			&id,
			&account_id,
			&total_price,
			&created_at,
			&product_id,
			&quantity,
		); err != nil {
			log.Println("ERROR: order repo GetOrdersByAccountID (rows.Scan): ", err)
			return nil, errors.New("error finding account's orders")
		}

		if order, exist := ordersMap[id]; !exist {
			ordersMap[id] = &Order{
				ID:         id,
				AccountID:  account_id,
				TotalPrice: total_price,
				CreatedAt:  created_at,
				Products: []OrderedProduct{
					{
						ID:       product_id,
						Quantity: quantity,
					},
				},
			}
		} else {
			order.Products = append(
				order.Products,
				OrderedProduct{
					ID:       product_id,
					Quantity: quantity,
				},
			)
		}
	}

	if err = rows.Err(); err != nil {
		log.Println("ERROR: order repo GetOrdersByAccountID (rows.Err): ", err)
		return nil, errors.New("error finding account's orders")
	}

	orders := []*Order{}

	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return orders, nil
}
