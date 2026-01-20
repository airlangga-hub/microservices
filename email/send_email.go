package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendEmail(msg amqp.Delivery) error {
	var order Order
	if err := json.Unmarshal(msg.Body, &order); err != nil {
		return fmt.Errorf("ERROR email SendEmail failed to unmarshal order: %w", err)
	}

	senderEmail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	recipientEmail := "customer@example.com"

	subject := "Subject: Payment Processed!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<h2>Order #%d Processed</h2><p>Total Price: %d</p>",
		order.ID, order.TotalPrice)

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", senderEmail, password, smtpServer)
	address := smtpServer + ":" + smtpPort

	err := smtp.SendMail(address, auth, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		return fmt.Errorf("ERROR email SendEmail failed to send email: %w", err)
	}

	return nil
}
