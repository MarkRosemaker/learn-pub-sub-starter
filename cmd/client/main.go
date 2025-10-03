package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connectionString = "amqp://guest:guest@localhost:5672/"
)

func main() {
	if err := do(); err != nil {
		log.Fatal(err)
	}
}

func do() error {
	log.Println("Starting Peril client...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return fmt.Errorf("dialing: %w", err)
	}
	defer conn.Close()

	log.Println("Successfully connected!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		return fmt.Errorf("welcoming client: %w", err)
	}

	log.Printf("client %q welcomed", username)

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	if _, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect,
		queueName, routing.PauseKey, pubsub.SimpleQueueTypeTransient,
	); err != nil {
		return err
	}

	log.Printf("queue %q declared and bound", queueName)

	util.WaitForSignal()

	log.Println("Shutting down and closing the connection...")

	return nil
}
