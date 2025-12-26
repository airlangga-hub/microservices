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

func NewClient() (*Client, error) {
	target := "localhost" + os.Getenv("ACCOUNT_PORT")

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		log.Fatalf("ERROR: account client NewClient: %v", err)
		return nil, errors.New("error creating grpc client connection")
	}

	service := pb.NewAccountServiceClient(conn)

	return &Client{
		Conn:    conn,
		Service: service,
	}, nil
}

func (c *Client) PostAccount(ctx context.Context, name string) (Account, error) {
	res, err := c.Service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		log.Println("ERROR: account client PostAccount: ", err)
		return Account{}, errors.New("error creating account")
	}

	return Account{ID: res.Account.Id, Name: res.Account.Name}, nil
}

func (c *Client) GetAccount(ctx context.Context, id int32) (Account, error) {
	res, err := c.Service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		log.Println("ERROR: account client GetAccount: ", err)
		return Account{}, errors.New("account not found")
	}
	
	return Account{ID: res.Account.Id, Name: res.Account.Name}, nil
}

func (c *Client) GetAccounts(ctx context.Context, offset, limit int32) ([]Account, error) {
	res, err := c.Service.GetAccounts(ctx, &pb.GetAccountsRequest{Offset: offset, Limit: limit})
	if err != nil {
		log.Println("ERROR: account client GetAccounts: ", err)
		return nil, errors.New("accounts not found")
	}
	
	accounts := []Account{}
	
	for _, account := range res.Accounts {
		accounts = append(accounts, Account{ID: account.Id, Name: account.Name})
	}
	
	return accounts, nil
}