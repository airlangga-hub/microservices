package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	accpb "github.com/airlangga-hub/microservices/order/account_pb"
	catpb "github.com/airlangga-hub/microservices/order/catalog_pb"
	"github.com/airlangga-hub/microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dbUrl := os.Getenv("ORDER_DB_URL")
	port := os.Getenv("ORDER_PORT")

	accountURL := os.Getenv("ACCOUNT_SERVICE_URL")
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")

	if accountURL == "" || catalogURL == "" || port == "" || dbUrl == "" {
		log.Fatalln("One or more order environment variables are missing")
	}

	repository, err := NewRepository(dbUrl)
	if err != nil {
		log.Fatalf("ERROR: order main: couldn't create repository: %v", err)
	}
	defer repository.Close()

	// account client
	accountConn, err := grpc.NewClient(accountURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: failed to create accountConn: ", err)
	}
	defer accountConn.Close()
	accountClient := accpb.NewAccountServiceClient(accountConn)

	// catalog client
	catalogConn, err := grpc.NewClient(catalogURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: failed to create catalogConn: ", err)
	}
	defer catalogConn.Close()
	catalogClient := catpb.NewCatalogServiceClient(catalogConn)

	// grpc server
	service := NewService(repository, accountClient, catalogClient)
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &Server{Svc: service})

	exitChan := make(chan error, 1)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		s := <-sig
		exitChan <- fmt.Errorf("received signal: %v", s)
	}()

	go func() {
		lis, _ := net.Listen("tcp", port)
		err := s.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			exitChan <- fmt.Errorf("grpc serve error: %v", err)
		}
	}()

	err = <-exitChan
	log.Printf("Shutting down. Reason: %v\n", err)

	s.GracefulStop()
}
