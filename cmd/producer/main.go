package main

import (
	"log"
	"time"

	"github.com/airlangga-hub/microservices/internal"
)

func main() {
	
	conn, err := internal.ConnectRabbit("angga", "Thebut12!", "localhost:5672", "customers")
	
	if err != nil {
		panic(err)
	}
	
	defer conn.Close()
	
	rabbitClient, err := internal.NewRabbitClient(conn)
	
	if err != nil {
		panic(err)
	}
	
	defer rabbitClient.Close()
	
	if err := rabbitClient.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}
	
	if err := rabbitClient.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}
	
	if err := rabbitClient.CreateBinding("customers_created", "customers.created.*", "customer_event"); err != nil {
		panic(err)
	}
	
	if err := rabbitClient.CreateBinding("customers_test", "customers.*", "customer_event"); err != nil {
		panic(err)
	}
	
	time.Sleep(time.Second * 10)
	
	log.Println(rabbitClient)
}