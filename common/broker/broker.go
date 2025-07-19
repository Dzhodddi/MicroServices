package broker

import (
	"commons/shared_errors"
	"context"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"time"
)

func New(addr string) (*rabbitmq.Connection, error) {
	conn, err := rabbitmq.Dial(addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type BrokerService struct {
	Conn *rabbitmq.Connection
}

func (b *BrokerService) Publish(ctx context.Context, msg []byte, queueName string) error {
	ch, err := b.Conn.Channel()
	if err != nil {
		return shared_errors.EstablishingConnection
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return shared_errors.QueueDeclareError
	}

	_ = ch.PublishWithContext(ctx, "", queue.Name, false, false, rabbitmq.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	return nil

}
