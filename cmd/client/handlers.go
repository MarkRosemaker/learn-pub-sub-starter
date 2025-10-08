package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ") // give user a new prompt
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")

		switch moveOutcome := gs.HandleMove(move); moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer, gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			if err := pubsub.PublishJSON(publishCh,
				routing.ExchangePerilTopic,
				util.DotJoin(routing.WarRecognitionsPrefix, gs.GetUsername()),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			); err != nil {
				log.Printf("ERROR: %v, requeuing %T", err, move)
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		default:
			log.Printf("ERROR: unknown move outcome %v", moveOutcome)
			return pubsub.NackDiscard
		}
	}
}

type GameLog struct {
	CurrentTime time.Time
	Message     string
	Username    string
}

func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(war gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")

		outcome, winner, looser := gs.HandleWar(war)
		msg := fmt.Sprintf("%s won a war against %s", winner, looser)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue // so another client can try to consume it
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeDraw:
			msg = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, looser)
			fallthrough
		case gamelogic.WarOutcomeOpponentWon, gamelogic.WarOutcomeYouWon:
			if err := pubsub.PublishGob(publishCh,
				routing.ExchangePerilTopic,
				util.DotJoin(routing.GameLogSlug, gs.GetUsername()),
				GameLog{
					CurrentTime: time.Now(),
					Message:     msg,
					Username:    gs.GetUsername(),
				},
			); err != nil {
				log.Printf("ERROR: %v, requeuing %T", err, war)
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		default:
			log.Printf("ERROR: unknown war outcome %v", outcome)
			return pubsub.NackDiscard
		}
	}
}
