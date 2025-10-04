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

func SubscribeJSON[T any](conn *amqp.Connection,
	exchange, queueName, key string, queueType SimpleQueueType, handler func(T),
) error {
	ch, q, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consuming %q: %w", q.Name, err)
	}

	go func() {
		for delivery := range deliveries {
			// Unmarshal the body of each message delivery into T
			var msg T
			if err := json.Unmarshal(delivery.Body, &msg); err != nil {
				panic(fmt.Errorf("unmarshaling delivery into %T: %w", msg, err))
			}

			// Handle the unmarshaled message
			handler(msg)

			// Acknowledge the message
			delivery.Ack(false)
		}
	}()

	return nil
}
