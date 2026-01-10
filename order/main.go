package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/airlangga-hub/microservices/order/pb"
	"google.golang.org/grpc"
)

func main() {
	dbUrl := os.Getenv("ORDER_DB_URL")
	port := os.Getenv("ORDER_PORT")

	repository, err := NewRepository(dbUrl)
	if err != nil {
		log.Fatalf("ERROR: order main: couldn't create repository: %v", err)
	}
	defer repository.Close()

	service := NewService(repository)

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
	log.Printf("Shutting down. Reason: %v", err)

	s.GracefulStop()
}
