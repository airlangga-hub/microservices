package catalog

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/airlangga-hub/microservices/services/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn    *grpc.ClientConn
	Service pb.CatalogServiceClient
}

func NewClient() (*Client, error) {
	target := "localhost" + os.Getenv("PORT")

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(target)
	if err != nil {
		log.Fatalf("ERROR: catalog client NewClient: %v", err)
		return nil, errors.New("error creating grpc client connection")
	}

	service := pb.NewCatalogServiceClient(conn)

	return &Client{
		Conn:    conn,
		Service: service,
	}, nil
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float32) (Product, error) {
	res, err := c.Service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		log.Println("ERROR: catalog client PostProduct: ", err)
		return Product{}, errors.New("error client post product")
	}

	return Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (Product, error) {
	res, err := c.Service.GetProduct(
		ctx,
		&pb.GetProductRequest{
			Id: id,
		},
	)
	if err != nil {
		log.Println("ERROR: catalog client GetProduct: ", err)
		return Product{}, errors.New("error client get product")
	}

	return Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, query string, ids []string, offset, limit int32) ([]Product, error) {
	res, err := c.Service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Offset: offset,
			Limit:  limit,
			Ids:    ids,
			Query:  query,
		},
	)
	if err != nil {
		log.Println("ERROR: catalog client GetProducts: ", err)
		return nil, errors.New("error client get products")
	}
	
	products := []Product{}
	
	for _, p := range res.Products {
		products = append(
			products,
			Product{
				ID: p.Id,
				Name: p.Name,
				Description: p.Description,
				Price: p.Price,
			},
		)
	}
	
	return products, nil
}
