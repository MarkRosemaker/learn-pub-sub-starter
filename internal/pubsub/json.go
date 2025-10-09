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

func SubscribeJSON[T any](conn *amqp.Connection,
	exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType,
) error {
	return subscribe(
		func(data []byte) (T, error) {
			var val T
			if err := json.Unmarshal(data, &val); err != nil {
				return val, err
			}
			return val, nil
		},
		conn, exchange, queueName, key, queueType, handler,
	)
}
