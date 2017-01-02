package strategy

import (
	"math/rand"

	"rummy"
	"rummy/deck"
)

type greedyStrategy struct {
	currentHand rummy.Hand
}

func newGreedyStrategy() Strategy {
	return &greedyStrategy{}
}

func (gs *greedyStrategy) PickUpCards(discardPile []deck.Card) int {
	return 0 // Always pick up from the stock.
}

func (gs *greedyStrategy) PlayCards(hand rummy.Hand) []deck.Card {
	// Play any melds in our hand.
	for _, m := range hand.Melds() {
		if len(hand) > len(m) {
			return m
		}
	}

	// Play any rummies in our hand.

	return nil
}

func (gs *greedyStrategy) Discard(hand rummy.Hand) deck.Card {
	// Discard a random card.
	i := rand.Intn(len(hand))
	toDiscard := hand.AsSlice()[i]
	delete(hand, toDiscard)
	gs.currentHand = hand
	return toDiscard
}
