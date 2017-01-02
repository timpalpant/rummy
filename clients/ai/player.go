package ai

import (
	"rummy"
	"rummy/clients/ai/strategy"
	"rummy/deck"
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
	events := make(chan *rummy.GameEvent, eventBufferSize)
	cp.g.Subscribe(events)

	for event := range events {
		if event.PlayerId == cp.playerId && event.Type == rummy.GameEvent_TURN_START {
			if err := cp.playTurn(); err != nil {
				return err
			}
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
	discardPile := valueSlice(cp.g.GameState().DiscardPile)
	n := cp.strategy.PickUpCards(discardPile)
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
	return cp.g.DiscardCard(cp.playerId, discard)
}
