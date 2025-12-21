package account

import (
	"log"
	"os"
)

func main() {
	dbUrl := os.Getenv("POSTGRES_DB_URL")
	port := os.Getenv("ACCOUNT_PORT")
	
	repository, err := NewRepository(dbUrl)
	if err != nil {
		log.Fatalf("ERROR: account main: couldn't create repository: %v", err)
	}
	defer repository.Close()
	
	service := NewService(repository)
	
	log.Fatal(ListenGrpc(service, port))
}