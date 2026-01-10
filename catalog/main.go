package catalog

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("CATALOG_PORT")

	repository, err := NewRepository()
	if err != nil {
		log.Fatalf("ERROR: catalog main: couldn't create repository: %v", err)
	}

	service := NewService(repository)
	
	log.Fatal(ListenGrpc(service, port))
}
