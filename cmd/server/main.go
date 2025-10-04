package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	log.Println("Starting Peril server...")

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return fmt.Errorf("dialing: %w", err)
	}
	defer conn.Close()

	log.Println("Successfully connected!")

	gamelogic.PrintServerHelp()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("getting channel: %w", err)
	}
	defer ch.Close()

	if err := gameLoop(ch); err != nil {
		return err
	}

	log.Println("Shutting down and closing the connection...")

	return nil
}

func gameLoop(ch *amqp.Channel) error {
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			log.Println("pausing the game")
			if err := pubsub.PublishJSON(ch,
				routing.ExchangePerilDirect, routing.PauseKey,
				routing.PlayingState{IsPaused: true}); err != nil {
				return fmt.Errorf("publishing pause state: %w", err)
			}
		case "resume":
			log.Println("resuming game")
			if err := pubsub.PublishJSON(ch,
				routing.ExchangePerilDirect, routing.PauseKey,
				routing.PlayingState{IsPaused: false}); err != nil {
				return fmt.Errorf("publishing resume state: %w", err)
			}
		case "quit":
			log.Println("exiting")
			return nil
		default:
			log.Printf("unknown command %q", words[0])
		}
	}
}
