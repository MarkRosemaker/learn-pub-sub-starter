package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

	gs := gamelogic.NewGameState(username)

	if err := pubsub.SubscribeJSON(conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, gs.GetUsername()),
		routing.PauseKey,
		pubsub.SimpleQueueTypeTransient, handlerPause(gs)); err != nil {
		return fmt.Errorf("subscribing to pause state: %w", err)
	}

	for {
		if err := gameLoop(gs); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		break
	}

	log.Println("Shutting down and closing the connection...")
	return nil
}

func gameLoop(gs *gamelogic.GameState) error {
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			if err := gs.CommandSpawn(words); err != nil {
				return fmt.Errorf("spawning: %w", err)
			}
		case "move":
			move, err := gs.CommandMove(words)
			if err != nil {
				return fmt.Errorf("moving: %w", err)
			}

			fmt.Printf("Move successful: %v\n", move)
		case "status":
			gs.CommandStatus()
		default:
			fmt.Printf("Command %q unknown\n", words[0])
			fallthrough
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return nil
		}
	}
}
