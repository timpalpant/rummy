package strategy

import (
	"math/rand"

	"github.com/timpalpant/rummy"
	"github.com/timpalpant/rummy/deck"
)

// Strategy is an interface for implementing Rummy 500 strategies.
type Strategy interface {
	// PickUpCards decides whether to:
	//   1) Take the top N cards from the discard pile, or
	//   2) Take the top card from the stock
	// If cards are taken from the discard pile, then the bottom-most card
	// (furthest into the discard pile) must be played in a meld this turn.
	//
	// PickUpCards should return the number of cards to pick up from the
	// discard pile. If 0, then a card will be picked up from teh stock.
	PickUpCards(discardPile []deck.Card) int

	// PlayCards selects, from the given Hand, cards to play.
	// The returned cards must either form a new Meld, or be playable
	// as rummies off of previously played Melds.
	PlayCards(hand rummy.Hand) []deck.Card

	// Discard selects a card from the hand to place in the discard pile.
	Discard(hand rummy.Hand) deck.Card

	// OnGameEvent is called for each GameEvent, and can be used in
	// advanced strategies to keep track of the cards that other players
	// pick up from the discard pile.
	//
	// OnGameEvent is only called for other player's actions. It is
	// assumed that the strategy knows its own actions.
	OnGameEvent(event *rummy.GameEvent)
}

// nopStrategy is a dummy strategy that picks up a card from the stock
// and then discards a random card from the hand.
type nopStrategy struct {
}

func newNopStrategy() Strategy {
	return &nopStrategy{}
}

func (ns *nopStrategy) PickUpCards(discardPile []deck.Card) int {
	return 0
}

func (ns *nopStrategy) PlayCards(hand rummy.Hand) []deck.Card {
	return nil
}

func (ns *nopStrategy) Discard(hand rummy.Hand) deck.Card {
	hs := hand.AsSlice()
	selected := rand.Intn(len(hs))
	return hs[selected]
}

func (ns *nopStrategy) OnGameEvent(event *rummy.GameEvent) {
}
