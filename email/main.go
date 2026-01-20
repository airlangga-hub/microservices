package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	emailQueueName := os.Getenv("EMAIL_QUEUE_NAME")
	amqpURL := os.Getenv("AMQP_URL")

	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")
	recipientEmail := os.Getenv("RECIPIENT_EMAIL")

	if emailQueueName == "" || amqpURL == "" || senderEmail == "" || senderPassword == "" || recipientEmail == "" {
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for msg := range messages {
			SendEmail(msg, senderEmail, senderPassword, recipientEmail)
		}
	}()

	sig := <-sigChan

	log.Printf("Signal %v received, exiting email service....\n", sig)
}
