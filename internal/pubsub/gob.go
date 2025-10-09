package pubsub

import (
	"bytes"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	return publish(
		func(in T) ([]byte, error) { return encodeGob(val) }, "application/gob",
		ch, exchange, key, val,
	)
}

func encodeGob[T any](val T) ([]byte, error) {
	w := &bytes.Buffer{}
	if err := gob.NewEncoder(w).Encode(val); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func decode[T any](data []byte) (T, error) {
	var val T
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&val); err != nil {
		return val, err
	}

	return val, nil
}
