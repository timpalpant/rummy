package deck

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
	return RankString(c.Rank) + SuitString(c.Suit)
}
