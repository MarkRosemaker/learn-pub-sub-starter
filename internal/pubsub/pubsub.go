package pubsub

import (
	"encoding/json/v2"
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

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
	q, err := ch.QueueDeclare(queueName, !isTransient, isTransient, isTransient, false,
		amqp.Table{
			// This will tell RabbitMQ to send failed messages to the dead letter exchange.
			"x-dead-letter-exchange": routing.ExchangePerilDL,
		})
	if err != nil {
		return ch, amqp.Queue{}, fmt.Errorf("declaring queue: %w", err)
	}

	if err := ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return ch, q, fmt.Errorf("binding queue: %w", err)
	}

	return ch, q, nil
}

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection,
	exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType,
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
			switch handler(msg) {
			case Ack: // Acknowledge the message
				log.Printf("acknowledging %T", msg)
				delivery.Ack(false)
			case NackRequeue:
				log.Printf("nack and requeue %T", msg)
				delivery.Nack(false, true)
			case NackDiscard:
				log.Printf("nack and discard %T", msg)
				delivery.Nack(false, false)
			}
		}
	}()

	return nil
}
