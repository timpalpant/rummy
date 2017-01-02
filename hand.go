package rummy

import (
	"fmt"
	"sort"

	"rummy/deck"
)

var nextRank = map[deck.Card_Rank]deck.Card_Rank{
	deck.Card_ACE:   deck.Card_TWO,
	deck.Card_TWO:   deck.Card_THREE,
	deck.Card_THREE: deck.Card_FOUR,
	deck.Card_FOUR:  deck.Card_FIVE,
	deck.Card_FIVE:  deck.Card_SIX,
	deck.Card_SIX:   deck.Card_SEVEN,
	deck.Card_SEVEN: deck.Card_EIGHT,
	deck.Card_EIGHT: deck.Card_NINE,
	deck.Card_NINE:  deck.Card_TEN,
	deck.Card_TEN:   deck.Card_JACK,
	deck.Card_JACK:  deck.Card_QUEEN,
	deck.Card_QUEEN: deck.Card_KING,
	deck.Card_KING:  deck.Card_ACE,
}

// Sort cards by suit and then rank.
type bySuitAndRank []deck.Card

func (b bySuitAndRank) Len() int {
	return len(b)
}

func (b bySuitAndRank) Less(i, j int) bool {
	if b[i].Suit == b[j].Suit {
		return b[i].Rank < b[j].Rank
	}
	return b[i].Suit < b[j].Suit
}

func (b bySuitAndRank) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Sort cards by rank only.
type byRank []deck.Card

func (b byRank) Len() int {
	return len(b)
}

func (b byRank) Less(i, j int) bool {
	return b[i].Rank < b[j].Rank
}

func (b byRank) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// TODO(palpant): Allow adjustable scoring rules, e.g. some stipulate that
// an Ace is worth 5 points if used as the low end of a run: A, 2, 3.
func value(card deck.Card) int {
	switch {
	case deck.Card_TWO <= card.Rank && card.Rank <= deck.Card_NINE:
		return 5
	case deck.Card_TEN <= card.Rank && card.Rank <= deck.Card_KING:
		return 10
	case card.Rank == deck.Card_ACE:
		return 15
	}

	// Shouldn't get here.
	return -1
}

// A Meld is a playable set of 3 or more cards.
// Melds must hold the cards in sorted order, to facilitate
// efficient checking of rummies.
type Meld []deck.Card

func (m Meld) Value() int {
	total := 0
	for _, card := range m {
		total += value(card)
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

	sort.Sort(bySuitAndRank(m))
	suit := m[0].Suit
	for i, card := range m[:len(m)-1] {
		if card.Suit != suit || nextRank[card.Rank] != m[i+1].Rank {
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
func (h Hand) Melds() []Meld {
	sets := h.Sets()
	runs := h.Runs()
	result := make([]Meld, 0, len(sets)+len(runs))
	result = append(result, sets...)
	result = append(result, runs...)
	return result
}

// Sets returns all possible sets in this Hand.
// The result may include overlapping sets of cards if there is a
// set of 4 cards of the same rank.
func (h Hand) Sets() []Meld {
	hs := h.AsSlice()
	sort.Sort(byRank(hs))

	result := make([]Meld, 0)
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
func (h Hand) Runs() []Meld {
	hs := h.AsSlice()
	sort.Sort(bySuitAndRank(hs))

	result := make([]Meld, 0)
	for i, card := range hs {
		potentialRun := []deck.Card{card}
		// Scan for a run of the same suit.
		for j := i + 1; j < len(hs); j++ {
			if hs[j].Suit == card.Suit && nextRank[hs[j-1].Rank] == hs[j].Rank {
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
	sort.Sort(bySuitAndRank(hs))
	cards := make([]string, len(hs))
	for i, c := range hs {
		cards[i] = deck.CardString(c)
	}

	return fmt.Sprintf("%v", cards)
}
