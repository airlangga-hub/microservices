package main

import (
	"log"
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	emailQueueName := os.Getenv("EMAIL_QUEUE_NAME")
	amqpURL := os.Getenv("AMQP_URL")
	if emailQueueName == "" || amqpURL == "" {
		log.Fatalln("FATAL email main environment variable missing")
	}
	
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalln("FATAL failed to connect amqp:", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalln("FATAL failed to create amqp channel:", err)
	}
	defer channel.Close()
	
	messages, err := channel.Consume(emailQueueName, "email-consumer", true, false, false, false, nil)
	if err != nil {
		log.Println("ERROR email main couldn't consume:", err)
		return
	}
}