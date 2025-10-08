package pubsub

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func publish[T any](
	marshal func(in T) ([]byte, error), contentType string,
	ch *amqp.Channel, exchange, key string, val T) error {
	b, err := marshal(val)
	if err != nil {
		return fmt.Errorf("marshaling %T: %w", val, err)
	}

	if err := ch.PublishWithContext(context.Background(),
		exchange, key, false, false,
		amqp.Publishing{ContentType: contentType, Body: b}); err != nil {
		return fmt.Errorf("publishing %T: %w", val, err)
	}

	return nil
}
