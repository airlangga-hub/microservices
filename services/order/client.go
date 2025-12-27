package order

import (
	"errors"
	"log"
	"os"

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
