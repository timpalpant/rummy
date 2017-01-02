package strategy

import (
	"math/rand"

	"rummy"
	"rummy/deck"
)

// greedyStrategy always plays melds and rummies in hand,
// and always picks up from the discard if any melds or rummies
// are possible.
type greedyStrategy struct {
	currentHand rummy.Hand
}

func newGreedyStrategy() Strategy {
	return &greedyStrategy{}
}

func (gs *greedyStrategy) PickUpCards(discardPile []deck.Card) int {
	for i := 1; i <= len(discardPile); i++ {
		bottomCard := len(discardPile) - i
		cardsToPickUp := discardPile[bottomCard:]
		hypotheticalHand := make(rummy.Hand, len(gs.currentHand)+i)
		for card := range gs.currentHand {
			hypotheticalHand[card] = struct{}{}
		}
		for _, card := range cardsToPickUp {
			hypotheticalHand[card] = struct{}{}
		}

		// Forms a meld, pick up from discard.
		melds := hypotheticalHand.Melds()
		if len(melds) > 0 {
			return i
		}
	}

	return 0 // Pick up from the stock.
}

// greedyStrategy always plays all possible melds and rummies.
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
	// Save current hand after discarding to assess picking up cards.
	gs.currentHand = hand
	return toDiscard
}

func (gs *greedyStrategy) OnGameEvent(event *rummy.GameEvent) {
}
