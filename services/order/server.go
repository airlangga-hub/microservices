package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/airlangga-hub/microservices/services/account"
	"github.com/airlangga-hub/microservices/services/catalog"
	"github.com/airlangga-hub/microservices/services/order/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	Svc           Service
	AccountClient *account.Client
	CatalogClient *catalog.Client
}

func ListenGrpc(service Service, port string) error {
	accountClient, err := account.NewClient()
	if err != nil {
		return fmt.Errorf("ERROR: order server ListenGrpc (account.NewClient): %v", err)
	}
	defer accountClient.Conn.Close()

	catalogClient, err := catalog.NewClient()
	if err != nil {
		return fmt.Errorf("ERROR: order server ListenGrpc (catalog.NewClient): %v", err)
	}
	defer catalogClient.Conn.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("ERROR: order server ListenGrpc (net.Listen): %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterOrderServiceServer(
		s,
		&Server{
			Svc:           service,
			AccountClient: accountClient,
			CatalogClient: catalogClient,
		},
	)

	return s.Serve(lis)
}

func (s *Server) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.AccountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}

	productIDs := []string{}
	mapIdQty := map[string]int32{}

	for _, p := range r.Products {
		productIDs = append(productIDs, p.Id)
		mapIdQty[p.Id] = p.Quantity
	}

	products, err := s.CatalogClient.GetProducts(
		ctx,
		"",
		productIDs,
		0,
		0,
	)
	if err != nil {
		return nil, err
	}

	orderedProducts := []OrderedProduct{}

	for _, p := range products {
		if qty, exist := mapIdQty[p.ID]; exist {
			orderedProducts = append(
				orderedProducts,
				OrderedProduct{
					ID: p.ID,
					Name: p.Name,
					Description: p.Description,
					Price: p.Price,
					Quantity: qty,
				},
			)
		}
	}

	order, err := s.Svc.PostOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		return nil, err
	}

	pbProducts := []*pb.OrderedProduct{}

	for _, p := range order.Products {
		pbProducts = append(
			pbProducts,
			&pb.OrderedProduct{
				Id: p.ID,
				Name: p.Name,
				Description: p.Description,
				Price: p.Price,
				Quantity: p.Quantity,
			},
		)
	}

	createdAt, err := order.CreatedAt.MarshalBinary()
	if err != nil {
		log.Println("ERROR: order server PostOrder (MarshalBinary): ", err)
		return nil, errors.New("error creating order")
	}

	return &pb.PostOrderResponse{
		Order: &pb.Order{
			Id: order.ID,
			AccountId: order.AccountID,
			Products: pbProducts,
			TotalPrice: order.TotalPrice,
			CreatedAt: createdAt,
		},
	}, nil
}

func (s *Server) GetOrdersByAccountID(ctx context.Context, r *pb.GetOrdersByAccountIDRequest) (*pb.GetOrdersByAccountIDResponse, error)
