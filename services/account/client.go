package account

import (
	"errors"
	"log"
	"os"

	"github.com/airlangga-hub/microservices/services/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn    *grpc.ClientConn
	Service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	port := os.Getenv("PORT")
	target := "localhost" + port

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		log.Fatalf("ERROR: account client NewClient: ", err)
		return nil, errors.New("error creating grpc client connection")
	}

	service := pb.NewAccountServiceClient(conn)

	return &Client{
		Conn:    conn,
		Service: service,
	}, nil
}
