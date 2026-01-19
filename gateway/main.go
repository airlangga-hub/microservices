package main

import (
	"log"
	"net/http"
	"os"

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

	// public endpoints
	http.HandleFunc("POST /api/signup", h.SignUp)
	http.HandleFunc("POST /api/login", h.Login)
	http.HandleFunc("GET /api/products", h.GetProducts)
	http.HandleFunc("GET /api/products/{id}", h.GetProductByID)
	http.HandleFunc("GET /api/products/search", h.SearchProducts)
	// buyer endpoints
	http.Handle("POST /api/order", h.AuthorizeMiddleware(http.HandlerFunc(h.CreateOrder)))
	http.Handle("GET /api/order", h.AuthorizeMiddleware(http.HandlerFunc(h.GetOrdersByAccountID)))
	// admin endpoints
	http.Handle("POST /admin/products", http.HandlerFunc(h.CreateProduct))

	log.Fatalln(http.ListenAndServe(port, nil))
}
