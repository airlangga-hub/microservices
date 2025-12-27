package order

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/airlangga-hub/microservices/services/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn    *grpc.ClientConn
	Service pb.OrderServiceClient
}

func NewClient() (*Client, error) {
	target := "localhost" + os.Getenv("ORDER_PORT")

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		log.Fatalf("ERROR: account client NewClient: %v", err)
		return nil, errors.New("error creating grpc client connection")
	}

	service := pb.NewOrderServiceClient(conn)

	return &Client{
		Conn:    conn,
		Service: service,
	}, nil
}

func (c *Client) PostOrder(ctx context.Context, accountID int32, products []OrderedProduct) (Order, error) {
	pbOrderedProducts := []*pb.OrderedProduct{}

	for _, p := range products {
		pbOrderedProducts = append(
			pbOrderedProducts,
			&pb.OrderedProduct{
				Id:       p.ID,
				Quantity: p.Quantity,
			},
		)
	}

	res, err := c.Service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountID,
			Products:  pbOrderedProducts,
		},
	)
	if err != nil {
		log.Println("ERROR: order client PostOrder (Service.PostOrder): ", err)
		return Order{}, errors.New("error creating order")
	}

	for i, pbOrderedProduct := range res.Order.Products {
		products[i] = OrderedProduct{
			ID:          pbOrderedProduct.Id,
			Name:        pbOrderedProduct.Name,
			Description: pbOrderedProduct.Description,
			Price:       pbOrderedProduct.Price,
			Quantity:    pbOrderedProduct.Quantity,
		}
	}

	t := time.Time{}
	if err := t.UnmarshalBinary(res.Order.CreatedAt); err != nil {
		log.Println("ERROR: order client PostOrder (t.UnmarshalBinary): ", err)
		return Order{}, errors.New("error creating order")
	}

	return Order{
		ID:         res.Order.Id,
		AccountID:  res.Order.AccountId,
		Products:   products,
		TotalPrice: res.Order.TotalPrice,
		CreatedAt:  t,
	}, nil
}

func (c *Client) GetOrdersByAccountID()
