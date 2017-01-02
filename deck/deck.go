package deck

import (
	"math/rand"
)

const (
	nSuits = 4
	nRanks = 13
)

type Deck []Card

func New() Deck {
	d := make(Deck, 0, nSuits*nRanks)
	for suit := 1; suit <= nSuits; suit++ {
		for rank := 1; rank <= nRanks; rank++ {
			card := Card{Card_Suit(suit), Card_Rank(rank)}
			d = append(d, card)
		}
	}
	return d
}

// Shuffle randomizes the order of the cards in the Deck
// using Fisher–Yates shuffle.
func (d Deck) Shuffle() {
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
}

// Take and return the top card from the Deck.
func (d *Deck) Pop() Card {
	old := *d
	n := len(old)
	top := old[n-1]
	*d = old[:n-1]
	return top
}
