package rummy

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/golang/glog"

	"rummy/deck"
	"rummy/meld"
)

const initialNumCards = 7

// Game manages the state machine for a single game of Rummy.
type Game struct {
	// The deck of cards that have not been picked up yet.
	stock deck.Deck
	// The discard pile.
	discard []deck.Card
	// The players in the Game.
	players []*player
	// map of player name -> id (index in players).
	name2id map[string]int32

	// Current player whose turn it is, and their turn state.
	turn                   int
	currentPlayer          int32
	currentPlayerTurnState GameState_TurnState
	// If the player picked up cards from the discard pile this turn,
	// then this is set to the bottom card that they picked up, since
	// this card must be played this turn before discarding.
	mustPlayCard *deck.Card
	// True once one player has "gone out" and the game is over.
	isOver bool

	// Subscribers to public game events.
	subscribers []chan *GameEvent
}

// NewGame initializes a new Game with a shuffled Deck of cards.
// There are initially no players. Players may join the game by
// calling AddPlayer, until the game is started by calling Deal.
func NewGame() *Game {
	d := deck.New()
	d.Shuffle()
	return &Game{
		stock:   d,
		name2id: make(map[string]int32),
		// No one can attempt to play until Deal is called.
		currentPlayer: -1,
	}
}

// AddPlayer adds a player with the given name to the game.
// AddPlayer can be called until Deal is called to add more players.
// Each player must have a unique name.
func (g *Game) AddPlayer(name string) (int32, error) {
	if g.currentPlayer != -1 {
		return 0, fmt.Errorf("game has already started, cannot join")
	}

	// Check if player with this name has joined already.
	if id, ok := g.name2id[name]; ok {
		return id, fmt.Errorf("player with name %v already joined", name)
	}

	p := &player{name: name}
	id := int32(len(g.players))
	g.name2id[name] = id
	g.players = append(g.players, p)
	return id, nil
}

// Deal starts the game, deals a hand to each player, and randomly
// selects the player to go first.
func (g *Game) Deal() error {
	if len(g.players) == 0 {
		return fmt.Errorf("no players in game")
	} else if len(g.players)*initialNumCards > len(g.stock) {
		return fmt.Errorf("too many players for deck: %v", len(g.players))
	}

	// Deal initial cards to each player.
	glog.Infof("Dealing %v cards to %v players", initialNumCards, len(g.players))
	for id, p := range g.players {
		p.hand = NewHand(g.stock[:initialNumCards])
		glog.Infof("Player %v (id: %v) initial hand: %s", p.name, id, p.hand)
		g.stock = g.stock[initialNumCards:]
	}

	// Initialize the discard pile.
	g.discard = []deck.Card{g.stock.Pop()}
	// Choose random player to start.
	g.currentPlayer = int32(rand.Intn(len(g.players)))
	// Notify any subscribers who goes first.
	// Anyone who subscribes after this will receive the event
	// upon subscription.
	g.publish(&GameEvent{
		PlayerId: g.currentPlayer,
		Type:     GameEvent_TURN_START,
	})
	return nil
}

// GameState returns the publicly observable state of the game.
func (g *Game) GameState() *GameState {
	playerStates := make([]*PlayerState, len(g.players))
	for i, p := range g.players {
		score := p.PublicScore()
		if g.isOver { // Final score is only revealed once game is over.
			score = p.Score()
		}
		playerStates[i] = &PlayerState{
			Id:             int32(i),
			Name:           p.name,
			Melds:          protoMelds(p.melds),
			Rummies:        protoSlice(p.rummies),
			NumCardsInHand: int32(len(p.hand)),
			CurrentScore:   int32(score),
		}
	}

	return &GameState{
		NumCardsInStock:   int32(len(g.stock)),
		DiscardPile:       protoSlice(g.discard),
		AggregatedMelds:   protoMelds(g.aggregatedMelds()),
		Players:           playerStates,
		Turn:              int32(g.turn),
		CurrentPlayerTurn: g.currentPlayer,
		TurnState:         g.currentPlayerTurnState,
		GameOver:          g.isOver,
	}
}

func protoMelds(melds []meld.Meld) []*Meld {
	result := make([]*Meld, len(melds))
	for i, m := range melds {
		result[i] = &Meld{
			Cards: protoSlice(m),
		}
	}
	return result
}

func (g *Game) PlayerHand(playerId int32) ([]deck.Card, error) {
	if playerId >= int32(len(g.players)) {
		return nil, fmt.Errorf("no such player: %v", playerId)
	}

	p := g.players[playerId]
	hs := p.hand.AsSlice()
	sort.Sort(deck.BySuitAndRank(hs))
	return hs, nil
}

func (g *Game) Subscribe(events chan *GameEvent) {
	g.subscribers = append(g.subscribers, events)
	if g.isOver {
		close(events)
	}
}

func (g *Game) publish(event *GameEvent) {
	for _, s := range g.subscribers {
		s <- event
	}
}

func (g *Game) PickUpStock(playerId int32) (deck.Card, error) {
	if g.currentPlayer != playerId {
		return deck.Card{}, fmt.Errorf("player %v, not %v turn", g.currentPlayer, playerId)
	}

	if g.currentPlayerTurnState != GameState_TURN_START {
		return deck.Card{}, fmt.Errorf("player %v has already picked up cards", playerId)
	}

	p := g.players[playerId]
	// Out of cards in the stock, shuffle the discard pile.
	if len(g.stock) == 0 {
		g.stock = deck.Deck(g.discard)
		g.stock.Shuffle()
		g.discard = []deck.Card{g.stock.Pop()}
	}

	card := g.stock.Pop()
	p.hand[card] = struct{}{}
	g.publish(&GameEvent{
		PlayerId: playerId,
		Type:     GameEvent_PICK_UP_STOCK,
		// Card is private to player, not published.
	})
	g.currentPlayerTurnState = GameState_PICKED_UP_CARDS
	return card, nil
}

func (g *Game) PickUpDiscard(playerId int32, nCards int) ([]deck.Card, error) {
	if g.currentPlayer != playerId {
		return nil, fmt.Errorf("player %v, not %v turn", g.currentPlayer, playerId)
	}

	if g.currentPlayerTurnState != GameState_TURN_START {
		return nil, fmt.Errorf("player %v has already picked up cards", playerId)
	}

	if nCards > len(g.discard) || nCards <= 0 {
		return nil, fmt.Errorf("can't pick up %v > %v cards in discard pile", nCards, len(g.discard))
	}

	p := g.players[playerId]
	lastCard := len(g.discard) - nCards
	cards := g.discard[lastCard:]
	// Save bottom card, which must be played this turn.
	mustPlayCard := cards[0]
	// Verify that player *can* play the mustPlayCard if cards are added to their hand.
	if !g.canPlayCard(mustPlayCard, p.hand, cards) {
		return nil, fmt.Errorf("player cannot play %v, the bottom card if %v "+
			"cards are picked up from the discard pile",
			deck.CardString(mustPlayCard), nCards)
	}

	// Pick up is valid, remove from discard pile and add to player's hand.
	g.discard = g.discard[:lastCard]
	for _, c := range cards {
		p.hand[c] = struct{}{}
	}
	// Save the card that must be played so we can verify that they play it before ending
	// their turn.
	g.mustPlayCard = &mustPlayCard

	g.publish(&GameEvent{
		PlayerId: playerId,
		Type:     GameEvent_PICK_UP_DISCARD,
		Cards:    protoSlice(cards),
	})
	g.currentPlayerTurnState = GameState_PICKED_UP_CARDS

	return cards, nil
}

func (g *Game) canPlayCard(card deck.Card, hand Hand, cards []deck.Card) bool {
	// Construct the hypothetical hand created by adding cards to hand.
	hypotheticalHand := make(Hand, len(hand)+len(cards))
	for card := range hand {
		hypotheticalHand[card] = struct{}{}
	}
	for _, card := range cards {
		hypotheticalHand[card] = struct{}{}
	}

	if _, ok := hypotheticalHand[card]; !ok {
		return false // Card is not in hand.
	}

	for _, meld := range hypotheticalHand.Melds() {
		for _, c := range meld {
			if card == c {
				return true // Card is playable as a meld.
			}
		}
	}

	// Card is not playable in any melds, check whether it can rummy
	// against any melds that have previously been played.
	return meld.CanRummy(card, g.aggregatedMelds())
}

func protoSlice(cards []deck.Card) []*deck.Card {
	result := make([]*deck.Card, len(cards))
	for i := range cards {
		result[i] = &cards[i]
	}
	return result
}

func (g *Game) PlayCards(playerId int32, cards []deck.Card) (int, error) {
	if g.currentPlayer != playerId {
		return 0, fmt.Errorf("player %v, not %v turn", g.currentPlayer, playerId)
	}

	if g.currentPlayerTurnState == GameState_TURN_START {
		return 0, fmt.Errorf("player %v must pick up cards before playing", playerId)
	}

	// Play must leave at least one card in hand for discard.
	p := g.players[playerId]
	if len(cards) >= len(p.hand) || len(cards) == 0 {
		return 0, fmt.Errorf("cannot play %d cards; hand contains %d",
			len(cards), len(p.hand))
	}

	// Validate that player has all cards they are trying to play,
	// and that all cards are unique.
	seenCards := make(map[deck.Card]struct{}, len(cards))
	for _, c := range cards {
		if _, ok := p.hand[c]; !ok {
			return 0, fmt.Errorf("player %v does not have %v in hand",
				playerId, deck.CardString(c))
		}

		if _, ok := seenCards[c]; ok {
			return 0, fmt.Errorf("invalid request: duplicate card %v", c)
		}
		seenCards[c] = struct{}{}
	}

	possibleMeld := meld.Meld(cards)
	isMeld := (possibleMeld.IsSet() || possibleMeld.IsRun())

	canRummy := true
	aggregatedMelds := g.aggregatedMelds()
	for _, c := range cards {
		canRummy = (canRummy && meld.CanRummy(c, aggregatedMelds))
	}

	if !isMeld && !canRummy {
		return 0, fmt.Errorf(
			"cannot play cards %v as a new meld or as rummies",
			ppCards(cards))
	}

	// At this point, this is a valid play. Remove the cards from the player's
	// hand and add them to played melds/rummies.
	for _, c := range cards {
		delete(p.hand, c)
		if g.mustPlayCard != nil && *g.mustPlayCard == c {
			g.mustPlayCard = nil // Card was played.
		}
	}

	if isMeld {
		p.melds = append(p.melds, possibleMeld)
	} else if canRummy {
		p.rummies = append(p.rummies, cards...)
	}
	score := possibleMeld.Value()
	g.publish(&GameEvent{
		PlayerId: playerId,
		Type:     GameEvent_PLAY_CARDS,
		Cards:    protoSlice(cards),
		Score:    int32(score),
	})
	g.currentPlayerTurnState = GameState_PLAYED_CARDS

	return score, nil
}

func ppCards(cards []deck.Card) string {
	result := make([]string, len(cards))
	for i, c := range cards {
		result[i] = deck.CardString(c)
	}

	return fmt.Sprintf("%v", result)
}

func (g *Game) DiscardCard(playerId int32, card deck.Card) error {
	if g.currentPlayer != playerId {
		return fmt.Errorf("player %v, not %v turn", g.currentPlayer, playerId)
	}

	// Discard is allowed only after cards have been picked up.
	if g.currentPlayerTurnState != GameState_PICKED_UP_CARDS &&
		g.currentPlayerTurnState != GameState_PLAYED_CARDS {
		return fmt.Errorf("player %v cannot discard", playerId)
	}

	if g.mustPlayCard != nil {
		return fmt.Errorf("player picked up and must play card %v before ending turn",
			deck.CardString(*g.mustPlayCard))
	}

	p := g.players[playerId]
	if _, ok := p.hand[card]; !ok {
		return fmt.Errorf("player %v cannot discard card %v not in hand",
			playerId, deck.CardString(card))
	}

	// Card is valid to discard, remove from hand and add it to the discard pile.
	delete(p.hand, card)
	g.discard = append(g.discard, card)
	g.publish(&GameEvent{
		PlayerId: playerId,
		Type:     GameEvent_DISCARD,
		Cards:    []*deck.Card{&card},
	})

	// If that was the last card in the player's hand, then the game is over.
	if len(p.hand) == 0 {
		g.endGame()
	} else {
		g.nextPlayer()
	}

	return nil
}

// Move Game forward to next player.
func (g *Game) nextPlayer() {
	g.turn++
	g.currentPlayer = (g.currentPlayer + 1) % int32(len(g.players))
	g.currentPlayerTurnState = GameState_TURN_START
	g.publish(&GameEvent{
		PlayerId: g.currentPlayer,
		Type:     GameEvent_TURN_START,
	})
}

func (g *Game) endGame() {
	g.isOver = true
	g.currentPlayer = -1
	g.publish(&GameEvent{
		Type: GameEvent_GAME_OVER,
	})
	for _, s := range g.subscribers {
		close(s)
	}
}

func (g *Game) CallRummy(playerId int32, cards []deck.Card) error {
	if playerId >= int32(len(g.players)) {
		return fmt.Errorf("no such player: %v", playerId)
	}

	// Verify that cards are in the discard pile, and that they are
	// playable off of an existing meld.
	discardSet := make(map[deck.Card]struct{}, len(g.discard))
	for _, c := range g.discard {
		discardSet[c] = struct{}{}
	}
	cardsSet := make(map[deck.Card]struct{}, len(cards))
	for _, c := range cards {
		if _, ok := discardSet[c]; !ok {
			return fmt.Errorf("card %v is not in the discard pile", deck.CardString(c))
		}
		cardsSet[c] = struct{}{}
	}

	// Verify that these cards can be played off of another meld.
	if !g.canPlay(cards) {
		return fmt.Errorf("cards: %v cannot be rummied", cards)
	}

	// Cards are a valid rummy, remove from the discard pile and
	// add them to the player's played cards.
	remainingDiscard := make([]deck.Card, 0, len(g.discard))
	for _, c := range g.discard {
		if _, ok := cardsSet[c]; !ok {
			remainingDiscard = append(remainingDiscard, c)
		}
	}
	g.discard = remainingDiscard
	p := g.players[playerId]
	p.rummies = append(p.rummies, cards...)
	return nil
}

func (g *Game) canPlay(cards []deck.Card) bool {
	possibleMeld := meld.Meld(cards)
	if possibleMeld.IsSet() || possibleMeld.IsRun() {
		return true
	}

	aggregatedMelds := g.aggregatedMelds()
	for _, c := range cards {
		if !meld.CanRummy(c, aggregatedMelds) {
			return false
		}
	}

	// All cards can be rummied.
	return true
}

// Get all of the extended melds in this Game, formed by taking the
// melds of each player and extending them with any rummies that have been
// played off of them.
func (g *Game) aggregatedMelds() []meld.Meld {
	melds := make([]meld.Meld, 0)
	rummies := make(map[deck.Card]struct{}, 0)
	for _, p := range g.players {
		for _, m := range p.melds {
			// Copy melds because we are going to add rummied cards to them
			// below and don't want to modify the meld attached to the player.
			copiedMeld := make(meld.Meld, len(m))
			copy(copiedMeld, m)
			melds = append(melds, copiedMeld)
		}
		for _, card := range p.rummies {
			rummies[card] = struct{}{}
		}
	}

	// Assign rummies to any matching sets.
	// TODO(palpant): It is possible for a rummied card to match both
	// a set and a run; this should be specified by the player when they
	// play the rummy. For now, we are assuming that rummies are against
	// the matching set in this case.
	for i, m := range melds {
		if m.IsSet() {
			// Find any matching cards with this rank.
			rank := m[0].Rank
			for card := range rummies {
				if card.Rank == rank {
					melds[i] = append(m, card)
					delete(rummies, card)
				}
			}
		}
	}

	// Assign remaining rummies to the runs they extend.
	for _, partialRun := range splitPartialRuns(rummies) {
		for i, m := range melds {
			if m.IsRun() {
				if deck.Sequential(m[len(m)-1], partialRun[0]) {
					// Extend run to the right.
					melds[i] = append(m, partialRun...)
					break
				} else if deck.Sequential(partialRun[len(partialRun)-1], m[0]) {
					// Extend run to the left.
					melds[i] = append(partialRun, m...)
					break
				}
			}
		}
	}

	return melds
}

// Takes the given set of rummy cards and splits them into contiguous
// partial runs, each having the same suit.
func splitPartialRuns(rummies map[deck.Card]struct{}) [][]deck.Card {
	if len(rummies) == 0 {
		return nil
	}

	runRummies := toSlice(rummies)
	sort.Sort(deck.BySuitAndRank(runRummies))

	result := make([][]deck.Card, 0)
	start := 0
	lastCard := runRummies[start]
	for i := 1; i < len(runRummies); i++ {
		card := runRummies[i]
		if !deck.Sequential(lastCard, card) {
			result = append(result, runRummies[start:i])
			start = i
		}
		lastCard = card
	}
	result = append(result, runRummies[start:])

	return result
}

func toSlice(cardSet map[deck.Card]struct{}) []deck.Card {
	s := make([]deck.Card, 0, len(cardSet))
	for card := range cardSet {
		s = append(s, card)
	}
	return s
}
