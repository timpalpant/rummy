package rummy

import (
	"rummy/deck"
	"rummy/meld"
	"rummy/scoring"
)

// player holds the state for a single player in a game of Rummy.
type player struct {
	name string
	// The cards in our hand. If len(hand) == 0, then the Game
	// is over.
	hand Hand
	// Played melds (sets or runs).
	melds []meld.Meld
	// Played rummies off of other melds.
	// The melds may be ones we have played, or ones that another
	// player in the Game has played.
	rummies []deck.Card
}

// Score returns the current score for this player, the sum of
// all played melds and rummies minus the score of cards still
// in hand.
func (p player) Score() int {
	total := p.PublicScore()
	for card := range p.hand {
		total -= scoring.Value(card)
	}

	return total
}

// PublicScore returns the publicly-visible score formed by
// taking the total of all played melds and rummies.
func (p player) PublicScore() int {
	total := 0
	for _, m := range p.melds {
		total += m.Value()
	}
	for _, card := range p.rummies {
		total += scoring.Value(card)
	}
	return total
}
