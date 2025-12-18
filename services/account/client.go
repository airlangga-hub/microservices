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
	GrpcConn   *grpc.ClientConn
	GrpcClient pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	port := os.Getenv("PORT")
	target := "localhost" + port
	
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	grpcConn, err := grpc.NewClient(target, opts...)
	if err != nil {
		log.Fatalf("ERROR: account client NewClient: ", err)
		return nil, errors.New("error creating grpc client connection")
	}
	
	grpcClient := pb.NewAccountServiceClient(grpcConn)
	
	return &Client{
		GrpcConn: grpcConn,
		GrpcClient: grpcClient,
	}, nil
}