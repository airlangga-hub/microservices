package order

import (
	"log"
	"os"
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
	
	log.Fatal(ListenGrpc(service, port))
}
