package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	accpb "github.com/airlangga-hub/microservices/gateway/account_pb"
	catpb "github.com/airlangga-hub/microservices/gateway/catalog_pb"
	orderpb "github.com/airlangga-hub/microservices/gateway/order_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	Secret        []byte
	AccountClient accpb.AccountServiceClient
	CatalogClient catpb.CatalogServiceClient
	OrderClient   orderpb.OrderServiceClient
}

func main() {
	port := os.Getenv("PORT")

	secret := os.Getenv("JWT_SECRET")
	jwtSecret := []byte(secret)

	accountURL := os.Getenv("ACCOUNT_SERVICE_URL")
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	orderURL := os.Getenv("ORDER_SERVICE_URL")

	if accountURL == "" || catalogURL == "" || orderURL == "" || secret == "" || port == "" {
		log.Fatalln("One or more gateway environment variables are missing")
	}

	// =====================
	// account client
	// =====================
	accountConn, err := grpc.NewClient(accountURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: failed to create accountConn: ", err)
	}
	defer accountConn.Close()
	accountClient := accpb.NewAccountServiceClient(accountConn)

	// =====================
	// catalog client
	// =====================
	catalogConn, err := grpc.NewClient(catalogURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: failed to create catalogConn: ", err)
	}
	defer catalogConn.Close()
	catalogClient := catpb.NewCatalogServiceClient(catalogConn)

	// =====================
	// order client
	// =====================
	orderConn, err := grpc.NewClient(orderURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: failed to create orderConn: ", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	// =====================
	// handler
	// =====================
	h := Handler{
		Secret:        jwtSecret,
		AccountClient: accountClient,
		CatalogClient: catalogClient,
		OrderClient:   orderClient,
	}

	// mux
	mux := http.NewServeMux()

	// public endpoints
	mux.HandleFunc("POST /api/signup", h.SignUp)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.HandleFunc("GET /api/products", h.GetProducts)
	mux.HandleFunc("GET /api/products/{id}", h.GetProductByID)
	mux.HandleFunc("GET /api/products/search", h.SearchProducts)
	// buyer endpoints
	mux.Handle("POST /api/order", h.AuthorizeMiddleware(http.HandlerFunc(h.CreateOrder)))
	mux.Handle("GET /api/order", h.AuthorizeMiddleware(http.HandlerFunc(h.GetOrdersByAccountID)))
	// admin endpoints
	mux.Handle("POST /admin/products", http.HandlerFunc(h.CreateProduct))

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// graceful shutdown
	exitChan := make(chan error, 1)

	go func() {
		log.Println("Listening on port: ", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			exitChan <- err
		}
	}()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		exitChan <- fmt.Errorf("Signal %v received, shutting down....\n", sig)
	}()

	err = <-exitChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Error shutting down server: ", err)
	}
	log.Println("Gracefully shutting down, error: ", err)
}
