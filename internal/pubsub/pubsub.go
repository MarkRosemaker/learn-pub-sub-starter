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
