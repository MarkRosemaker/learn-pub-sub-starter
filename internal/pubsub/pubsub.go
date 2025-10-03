package pubsub

import (
	"context"
	"encoding/json/v2"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	b, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("marshaling %T: %w", val, err)
	}

	if err := ch.PublishWithContext(context.Background(),
		exchange, key, false, false,
		amqp.Publishing{ContentType: "application/json", Body: b}); err != nil {
		return fmt.Errorf("publishing %T: %w", val, err)
	}

	return nil
}

// an enum to represent "durable" or "transient"
type SimpleQueueType int

const (
	SimpleQueueTypeDurable SimpleQueueType = iota
	SimpleQueueTypeTransient
)

func DeclareAndBind(conn *amqp.Connection,
	exchange, queueName, key string, queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("getting channel: %w", err)
	}
	defer ch.Close()

	isTransient := queueType == SimpleQueueTypeTransient
	q, err := ch.QueueDeclare(queueName, !isTransient, isTransient, isTransient, false, nil)
	if err != nil {
		return ch, amqp.Queue{}, fmt.Errorf("declaring queue: %w", err)
	}

	if err := ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return ch, q, fmt.Errorf("binding queue: %w", err)
	}

	return ch, q, nil
}
