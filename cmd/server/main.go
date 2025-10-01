package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Wait for a signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	fmt.Println("\nShutting down and closing the connection...")

	return nil
}
