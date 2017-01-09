package rummy

import (
	"fmt"
	"sort"

	"github.com/timpalpant/rummy/deck"
	"github.com/timpalpant/rummy/meld"
)

type Hand map[deck.Card]struct{}

func NewHand(cards []deck.Card) Hand {
	h := make(Hand, len(cards))
	for _, c := range cards {
		h[c] = struct{}{}
	}
	return h
}

// Return a copy of this Hand as a slice of cards.
func (h Hand) AsSlice() []deck.Card {
	s := make([]deck.Card, 0, len(h))
	for card := range h {
		s = append(s, card)
	}
	return s
}

// Melds returns all possible sets and runs in this Hand.
// The resulting melds may include overlapping sets of cards.
func (h Hand) Melds() []meld.Meld {
	sets := h.Sets()
	runs := h.Runs()
	result := make([]meld.Meld, 0, len(sets)+len(runs))
	result = append(result, sets...)
	result = append(result, runs...)
	return result
}

// Sets returns all possible sets in this Hand.
// The result may include overlapping sets of cards if there is a
// set of 4 cards of the same rank.
func (h Hand) Sets() []meld.Meld {
	hs := h.AsSlice()
	sort.Sort(deck.ByRank(hs))

	result := make([]meld.Meld, 0)
	for i, card := range hs {
		potentialSet := []deck.Card{card}
		// Scan for other cards of the same rank.
		for j := i + 1; j < len(hs); j++ {
			if hs[j].Rank == card.Rank {
				potentialSet = append(potentialSet, hs[j])
			} else {
				break // Cards are sorted by Rank.
			}
		}

		if len(potentialSet) >= 3 {
			result = append(result, potentialSet)
		}
	}
	return result
}

// Runs returns all possible runs in this Hand.
// The result may include overlapping sets of Cards if there is a run
// of > 3 cards.
func (h Hand) Runs() []meld.Meld {
	hs := h.AsSlice()
	sort.Sort(deck.BySuitAndRank(hs))

	result := make([]meld.Meld, 0)
	for i, card := range hs {
		potentialRun := []deck.Card{card}
		// Scan for a run of the same suit.
		for j := i + 1; j < len(hs); j++ {
			if deck.Sequential(hs[j-1], hs[j]) {
				potentialRun = append(potentialRun, hs[j])
			} else {
				break // Cards are sorted by suit and rank.
			}
		}

		if len(potentialRun) >= 3 {
			result = append(result, potentialRun)
		}
	}
	return result
}

func (h Hand) String() string {
	hs := h.AsSlice()
	sort.Sort(deck.BySuitAndRank(hs))
	cards := make([]string, len(hs))
	for i, c := range hs {
		cards[i] = deck.CardString(c)
	}

	return fmt.Sprintf("%v", cards)
}
