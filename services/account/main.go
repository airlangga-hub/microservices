package account

import (
	"log"
	"os"
)

func main() {
	dbUrl := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	
	repository, err := NewRepository(dbUrl)
	if err != nil {
		log.Fatalf("ERROR: account main: couldn't create repository")
	}
	defer repository.Close()
	
	service := NewService(repository)
	
	log.Fatal(ListenGrpc(service, port))
}