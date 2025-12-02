package main

import (
	"log"
	"os"

	"github.com/airlangga-hub/microservices/internal"
	"github.com/joho/godotenv"
)

func main() {
	
	godotenv.Load()
	
	password := os.Getenv("PASSWORD")
	username := os.Getenv("USERNAME")
	host := os.Getenv("HOST")
	vhost := os.Getenv("VHOST")
	
	conn, err := internal.ConnectRabbit(username, password, host, vhost)
	
	if err != nil {
		panic(err)
	}
	
	defer conn.Close()
	
	rabbitClient, err := internal.NewRabbitClient(conn)
	
	if err != nil {
		panic(err)
	}
	
	defer rabbitClient.Close()
	
	messageBus, err := rabbitClient.Consume("customers_created", "email_service", false)
	
	if err != nil {
		panic(err)
	}
	
	blocking := make(chan bool)
	
	go func() {
		for message := range messageBus {
			log.Println("New message: ", message)
			
			if err := message.Ack(false); err != nil {
				log.Println("Acknowledge message failed: ", err)
				continue
			}
			
			log.Println("Acknowledged message: ", message.MessageId)
		}
	}()

	<-blocking
}