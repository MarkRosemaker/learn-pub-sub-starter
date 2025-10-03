package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	if err := do(); err != nil {
		log.Fatal(err)
	}
}

func do() error {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return fmt.Errorf("dialing: %w", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected!")

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("getting channel: %w", err)
	}
	defer ch.Close()

	if err := pubsub.PublishJSON(ch,
		routing.ExchangePerilDirect, routing.PauseKey,
		routing.PlayingState{IsPaused: true}); err != nil {
		return fmt.Errorf("publishing initial pause state: %w", err)
	}

	// Wait for a signal to exit
	fmt.Println("Press Ctrl+C to exit")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	fmt.Println("\nShutting down and closing the connection...")

	return nil
}
