package pubsub

import (
	"encoding/json/v2"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	return publish(
		func(in T) ([]byte, error) { return json.Marshal(val) }, "application/json",
		ch, exchange, key, val,
	)
}
