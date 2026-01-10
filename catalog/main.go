package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/airlangga-hub/microservices/catalog/pb"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("CATALOG_PORT")

	repository, err := NewRepository()
	if err != nil {
		log.Fatalf("ERROR: catalog main: couldn't create repository: %v", err)
	}

	service := NewService(repository)

	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, &Server{Svc: service})

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
