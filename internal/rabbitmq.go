package internal

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn		*amqp.Connection
	chann		*amqp.Channel
}

func ConnectRabbit(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitClient(conn *amqp.Connection) (RabbitClient, error) {
	
	chann, err := conn.Channel()
	
	if err != nil {
		return RabbitClient{}, err
	}
	
	return RabbitClient{
		conn: conn,
		chann: chann,
	}, nil
}

func (r RabbitClient) Close() error {
	return r.chann.Close()
}

func (r RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := r.chann.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	return err
}

func (r RabbitClient) CreateBinding(name, key, exchange string) error {
	return r.chann.QueueBind(name, key, exchange, false, nil)
}

func (r RabbitClient) Send(ctx context.Context, exchange, key string, options amqp.Publishing) error {
	return r.chann.PublishWithContext(
		ctx,
		exchange,
		key,
		true,
		false,
		options,
	)
}

func (r RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.chann.Consume(queue, consumer, autoAck, false, false, false, nil)
}