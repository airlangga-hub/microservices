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
	
	emailQueueName := os.Getenv("EMAIL_QUEUE_NAME")
	amqpURL := os.Getenv("AMQP_URL")

	if accountURL == "" || catalogURL == "" || port == "" || dbUrl == "" || emailQueueName == "" || amqpURL == "" {
		log.Fatalln("One or more order environment variables are missing")
	}
	
	// publisher amqp
	publisher, err := NewPublisher(amqpURL, emailQueueName)
	if err != nil {
		log.Fatalln("FATAL order main couldn't create publisher:", err)
	}
	defer publisher.Connection.Close()
	// declare channel
	_, err = publisher.Channel.QueueDeclare(emailQueueName, true, false, false, false, nil)
	if err != nil {
		log.Println("FATAL order main couldn't declare queue:", err)
		return
	}
	
	// repository
	repository, err := NewRepository(dbUrl)
	if err != nil {
		log.Fatalln("FATAL order main couldn't create repository:", err)
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
	service := NewService(repository, accountClient, catalogClient, publisher)
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &Server{Svc: service})

	// graceful exit
	exitChan := make(chan error, 1)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		s := <-sig
		exitChan <- fmt.Errorf("received signal: %v", s)
	}()

	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			exitChan <- fmt.Errorf("error tcp listen: %v", err)
		}

		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			exitChan <- fmt.Errorf("grpc serve error: %v", err)
		}
	}()

	err = <-exitChan
	log.Printf("Shutting down. Reason: %v\n", err)

	s.GracefulStop()
}
