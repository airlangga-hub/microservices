package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Connection     *amqp.Connection
	Channel        *amqp.Channel
	EmailQueueName string
}

func NewPublisher(url, emailQueueName string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println("FATAL failed to connect amqp:", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Println("FATAL failed to create amqp channel:", err)
		return nil, err
	}

	return &Publisher{
		Connection:     conn,
		Channel:        channel,
		EmailQueueName: emailQueueName,
	}, nil
}

func (p *Publisher) Publish(message any) {
	marshaled, err := json.Marshal(message)
	if err != nil {
		log.Println("ERROR order publisher Publish (Marshal):", err)
	}

	if err := p.Channel.Publish(
		"",
		p.EmailQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshaled,
		},
	); err != nil {
		log.Println("ERROR order publisher Publish (Publish):", err)
	}
}
