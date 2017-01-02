package strategy

import (
	"math/rand"

	"rummy"
	"rummy/deck"
)

// Strategy is an interface for implementing Rummy 500 strategies.
type Strategy interface {
	// PickUpCards either:
	//   1) Takes the top N cards from the discard pile
	//   2) Takes the top card from the stock, or
	// If cards are taken from the discard pile, then the bottom-most card
	// (furthest into the discard pile) must be played in a meld this turn.
	// If no cards are picked up from the discard pile, then one will
	// automatically be picked up from the stock.
	PickUpCards(discardPile []deck.Card) int

	// PlayCards selects, from the given Hand, the
	// cards to play as points this turn. Note that the played cards
	// must not be overlapping, but may be rummies off of previously plays.
	PlayCards(hand rummy.Hand) []deck.Card

	// Discard selects a card from the hand to place in the discard pile.
	Discard(hand rummy.Hand) deck.Card
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
