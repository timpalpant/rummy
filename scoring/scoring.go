package scoring

import (
	"github.com/timpalpant/rummy/deck"
)

// TODO(palpant): Allow adjustable scoring rules, e.g. some stipulate that
// an Ace is worth 5 points if used as the low end of a run: A, 2, 3.
func Value(card deck.Card) int {
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
