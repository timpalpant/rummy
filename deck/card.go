package deck

import (
	"github.com/fatih/color"
)

var nextRank = map[Card_Rank]Card_Rank{
	Card_ACE:   Card_TWO,
	Card_TWO:   Card_THREE,
	Card_THREE: Card_FOUR,
	Card_FOUR:  Card_FIVE,
	Card_FIVE:  Card_SIX,
	Card_SIX:   Card_SEVEN,
	Card_SEVEN: Card_EIGHT,
	Card_EIGHT: Card_NINE,
	Card_NINE:  Card_TEN,
	Card_TEN:   Card_JACK,
	Card_JACK:  Card_QUEEN,
	Card_QUEEN: Card_KING,
	Card_KING:  Card_ACE,
}

func NextRank(rank Card_Rank) Card_Rank {
	return nextRank[rank]
}

func Sequential(card1, card2 Card) bool {
	return card1.Suit == card2.Suit && NextRank(card1.Rank) == card2.Rank
}

var rankToStr = [...]string{
	"X",
	"A",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"10",
	"J",
	"Q",
	"K",
}

var suitToStr = [...]string{
	"X",
	"♥",
	"♦",
	"♣",
	"♠",
}

func SuitString(s Card_Suit) string {
	return suitToStr[s]
}

func RankString(r Card_Rank) string {
	return rankToStr[r]
}

func CardString(c Card) string {
	s := RankString(c.Rank) + SuitString(c.Suit)
	if c.Suit == Card_HEARTS || c.Suit == Card_DIAMONDS {
		return color.RedString(s)
	}
	return s
}
