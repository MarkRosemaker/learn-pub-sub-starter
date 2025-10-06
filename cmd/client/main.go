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

	publishCh, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("getting channel: %w", err)
	}
	defer publishCh.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		return fmt.Errorf("welcoming client: %w", err)
	}

	gs := gamelogic.NewGameState(username)

	if err := pubsub.SubscribeJSON(conn,
		routing.ExchangePerilTopic,                              // exchange
		util.DotJoin(routing.ArmyMovesPrefix, gs.GetUsername()), // queue name
		util.DotJoin(routing.ArmyMovesPrefix, util.Asterisk),    // key
		pubsub.SimpleQueueTypeTransient,
		handlerMove(gs, publishCh),
	); err != nil {
		return fmt.Errorf("subscribing to moves: %w", err)
	}

	if err := pubsub.SubscribeJSON(conn,
		routing.ExchangePerilTopic,                                 // exchange
		routing.WarRecognitionsPrefix,                              // queue name
		util.DotJoin(routing.WarRecognitionsPrefix, util.Asterisk), // key
		pubsub.SimpleQueueTypeDurable,
		handlerRecognitionOfWar(gs),
	); err != nil {
		return fmt.Errorf("subscribing to war recognitions: %w", err)
	}

	if err := pubsub.SubscribeJSON(conn,
		routing.ExchangePerilDirect,                      // exchange
		util.DotJoin(routing.PauseKey, gs.GetUsername()), // queue name
		routing.PauseKey,                                 // key
		pubsub.SimpleQueueTypeTransient,
		handlerPause(gs),
	); err != nil {
		return fmt.Errorf("subscribing to pauses: %w", err)
	}

	for {
		if err := gameLoop(publishCh, gs); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		break
	}

	log.Println("Shutting down and closing the connection...")
	return nil
}

func gameLoop(ch *amqp.Channel, gs *gamelogic.GameState) error {
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
			mv, err := gs.CommandMove(words)
			if err != nil {
				return fmt.Errorf("moving: %w", err)
			}

			if err := pubsub.PublishJSON(ch,
				routing.ExchangePerilTopic,
				util.DotJoin(routing.ArmyMovesPrefix, gs.GetUsername()),
				mv,
			); err != nil {
				return fmt.Errorf("publishing move: %w", err)
			}

			fmt.Printf("Moved %v units to %s\n", len(mv.Units), mv.ToLocation)
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
