package meld

import (
	"fmt"
	"sort"

	"rummy/deck"
	"rummy/scoring"
)

// A Meld is a playable set of 3 or more cards.
// Melds must hold the cards in sorted order, to facilitate
// efficient checking of rummies.
type Meld []deck.Card

func (m Meld) Value() int {
	total := 0
	for _, card := range m {
		total += scoring.Value(card)
	}

	return total
}

func (m Meld) IsSet() bool {
	if len(m) < 3 {
		return false
	}

	rank := m[0].Rank
	for _, card := range m {
		if card.Rank != rank {
			return false
		}
	}

	return true
}

func (m Meld) IsRun() bool {
	if len(m) < 3 {
		return false
	}

	sort.Sort(deck.BySuitAndRank(m))
	for i, card := range m[:len(m)-1] {
		if !deck.Sequential(card, m[i+1]) {
			return false
		}
	}

	return true
}

func (m Meld) String() string {
	cards := make([]string, len(m))
	for i, c := range m {
		cards[i] = deck.CardString(c)
	}

	return fmt.Sprintf("%v", cards)
}

// Check whether the given card can be played as a rummy off of
// any of the given melds.
func CanRummy(card deck.Card, melds []Meld) bool {
	for _, m := range melds {
		if m.IsSet() && m[0].Rank == card.Rank {
			return true
		} else if deck.Sequential(card, m[0]) || deck.Sequential(m[len(m)-1], card) {
			return true
		}
	}

	return false

}
