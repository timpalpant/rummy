package ai

import (
	"github.com/timpalpant/rummy"
	"github.com/timpalpant/rummy/clients/ai/strategy"
	"github.com/timpalpant/rummy/deck"

	"github.com/golang/glog"
)

// Size of channel buffer for game events.
// This must be large enough to buffer all of the events that
// a CP takes during a single turn, otherwise the engine will
// lock up because the player does not drain their own events
// while playing.
const eventBufferSize = 1000

// Play the given game with the given strategy.
func PlayGame(g *rummy.Game, playerId int32, strategy strategy.Strategy) error {
	p := &computerPlayer{g, playerId, strategy}
	return p.Play()
}

// computerPlayer automatically initiates gameplay actions when it
// is their turn according to a certain strategy.
type computerPlayer struct {
	g        *rummy.Game
	playerId int32
	strategy strategy.Strategy
}

func (cp *computerPlayer) Play() error {
	// TODO(palpant): Fix race, computer needs to be subscribed to events
	// before game is dealt, otherwise this can deadlock because the computer
	// will miss the event that signals it to start its turn.
	events := make(chan *rummy.GameEvent, eventBufferSize)
	cp.g.Subscribe(events)

	for event := range events {
		if event.PlayerId == cp.playerId {
			if event.Type == rummy.GameEvent_TURN_START {
				if err := cp.playTurn(); err != nil {
					return err
				}
			}
		} else {
			cp.strategy.OnGameEvent(event)
		}
	}

	return nil
}

func valueSlice(cards []*deck.Card) []deck.Card {
	result := make([]deck.Card, len(cards))
	for i, c := range cards {
		result[i] = *c
	}
	return result
}

func (cp *computerPlayer) playTurn() error {
	glog.V(1).Infof("Starting CP turn")
	discardPile := valueSlice(cp.g.GameState().DiscardPile)
	n := cp.strategy.PickUpCards(discardPile)
	glog.V(1).Infof("CP chose to pick up %d cards from discard", n)
	if n > 0 {
		_, err := cp.g.PickUpDiscard(cp.playerId, n)
		if err != nil {
			return err
		}
	} else {
		_, err := cp.g.PickUpStock(cp.playerId)
		if err != nil {
			return err
		}
	}

	for {
		hand, err := cp.g.PlayerHand(cp.playerId)
		if err != nil {
			return err
		}
		cards := cp.strategy.PlayCards(rummy.NewHand(hand))
		if len(cards) == 0 {
			break
		}

		glog.V(1).Infof("CP chose to play cards: %v", cards)
		_, err = cp.g.PlayCards(cp.playerId, cards)
		if err != nil {
			return err
		}
	}

	hand, err := cp.g.PlayerHand(cp.playerId)
	if err != nil {
		return err
	}
	discard := cp.strategy.Discard(rummy.NewHand(hand))
	glog.V(1).Infof("CP chose to dicard: %v", deck.CardString(discard))
	return cp.g.DiscardCard(cp.playerId, discard)
}
