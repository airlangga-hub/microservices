package main

import (
	"log"
	"net/http"
	"os"

	accpb "github.com/airlangga-hub/microservices/gateway/account_pb"
	catpb "github.com/airlangga-hub/microservices/gateway/catalog_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Secret        []byte
	AccountClient accpb.AccountServiceClient
	CatalogClient catpb.CatalogServiceClient
}

func main() {
	secret := os.Getenv("JWT_SECRET")
	jwtSecret := []byte(secret)

	accountConn, err := grpc.NewClient("account:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("FATAL: failed to create accountConn")
	}
	defer accountConn.Close()

	accountClient := accpb.NewAccountServiceClient(accountConn)

	catalogConn, err := grpc.NewClient("account:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("FATAL: failed to create accountConn")
	}
	defer catalogConn.Close()

	catalogClient := catpb.NewCatalogServiceClient(catalogConn)

	cfg := Config{
		Secret:        jwtSecret,
		AccountClient: accountClient,
		CatalogClient: catalogClient,
	}

	// public endpoints
	http.HandleFunc("POST /api/signup", cfg.SignUp)
	http.HandleFunc("POST /api/login", cfg.Login)
	http.HandleFunc("GET /api/products", cfg.GetProducts)
	// buyer endpoints
	// seller endpoints
}
