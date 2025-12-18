package account

import (
	"context"
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

func (c *Client) PostAccount(ctx context.Context, name string) (*pb.Account, error) {
	res, err := c.Service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		log.Println("ERROR: account client PostAccount: ", err)
		return nil, errors.New("error client post account")
	}

	return res.Account, nil
}

func (c *Client) GetAccount(ctx context.Context, id int32) (*pb.Account, error) {
	res, err := c.Service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		log.Println("ERROR: account client GetAccount: ", err)
		return nil, errors.New("error client get account")
	}
	
	return res.Account, nil
}

func (c *Client) GetAccounts(ctx context.Context, offset, limit int32) ([]*pb.Account, error) {
	res, err := c.Service.GetAccounts(ctx, &pb.GetAccountsRequest{Offset: offset, Limit: limit})
	if err != nil {
		log.Println("ERROR: account client GetAccounts: ", err)
		return nil, errors.New("error client get accounts")
	}
	
	return res.Accounts, nil
}