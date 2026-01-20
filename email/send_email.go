package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendEmail(msg amqp.Delivery, senderEmail, senderPassword, recipientEmail string) error {
	var order Order
	if err := json.Unmarshal(msg.Body, &order); err != nil {
		return fmt.Errorf("ERROR email SendEmail failed to unmarshal order: %w", err)
	}

	var itemsList string
	for _, p := range order.Products {
		itemsList += fmt.Sprintf("- %s (Qty: %d) - $%.2f\n", p.Name, p.Quantity, float64(p.Price)/100)
	}

	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Order Confirmation #" + fmt.Sprint(order.ID) + "\n"
	body := fmt.Sprintf(
		"Thank you for your order!\n\n"+
			"Items:\n%s\n"+
			"Total: $%.2f\n",
		itemsList,
		float64(order.TotalPrice)/100,
	)

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	address := smtpServer + ":" + smtpPort

	err := smtp.SendMail(address, auth, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		return fmt.Errorf("ERROR email SendEmail failed to send email: %w", err)
	}

	return nil
}
