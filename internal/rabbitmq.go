package internal

import (
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

func (r RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := r.chann.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

func (r RabbitClient) CreateBinding(name, binding, exchange string) error {
	return r.chann.QueueBind(name, binding, exchange, false, nil)
}