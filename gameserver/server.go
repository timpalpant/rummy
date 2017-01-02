package gameserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"rummy"
	"rummy/clients/ai"
	"rummy/clients/ai/strategy"
	"rummy/deck"
)

const (
	gameExpiration   = time.Hour
	eventsBufferSize = 100
)

// RummyServer implements RummyService.
// RummyServer creates new Games of rummy and provides an interface
// for players to perform actions on the Game.
type RummyServer struct {
	connStr string

	// gamesMu protects games and completedGames.
	gamesMu sync.Mutex
	// map of game name -> Game.
	games map[string]*rummy.Game
	// map of game name -> expiration time at which it will be deleted.
	completedGames map[string]time.Time
}

func NewRummyServer(connStr string) *RummyServer {
	return &RummyServer{
		connStr:        connStr,
		games:          make(map[string]*rummy.Game),
		completedGames: make(map[string]time.Time),
	}
}

func (s *RummyServer) CreateGame(ctx context.Context, req *rummy.CreateGameRequest) (*rummy.CreateGameResponse, error) {
	glog.V(1).Infof("CreateGame: %v", req)

	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	s.garbageCollectCompletedGames()

	if _, ok := s.games[req.GameName]; ok {
		return nil, fmt.Errorf("game %v already exists", req.GameName)
	}

	s.games[req.GameName] = rummy.NewGame()
	return &rummy.CreateGameResponse{}, nil
}

// Must be called while holding gamesMu.
func (s *RummyServer) garbageCollectCompletedGames() {
	for name, game := range s.games {
		if game.GameState().GameOver {
			glog.Infof("Detected that game %v is over", name)
			s.completedGames[name] = time.Now()
		}
	}

	now := time.Now()
	for name, expiration := range s.completedGames {
		glog.V(1).Infof("Completed game %v expires at %v", name, expiration)
		if expiration.After(now) {
			glog.Infof("Removing completed game %v", name)
			delete(s.games, name)
			delete(s.completedGames, name)
		}
	}
}

func (s *RummyServer) JoinGame(ctx context.Context, req *rummy.JoinGameRequest) (*rummy.JoinGameResponse, error) {
	glog.V(1).Infof("JoinGame: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	id, err := g.AddPlayer(req.PlayerName)
	if err == nil && req.Strategy != "" {
		// Start computer player.
		glog.Infof("Starting computer player %v for game %v with strategy %v",
			req.PlayerName, req.GameName, req.Strategy)
		strat, err := strategy.ForName(req.Strategy)
		if err != nil {
			return nil, err
		}

		go ai.PlayGame(s.connStr, req.GameName, id, strat)
	}

	return &rummy.JoinGameResponse{
		PlayerId: id,
	}, err
}

func (s *RummyServer) StartGame(ctx context.Context, req *rummy.StartGameRequest) (*rummy.StartGameResponse, error) {
	glog.V(1).Infof("StartGame: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	glog.Infof("Starting game: %v", req.GameName)
	err := g.Deal()
	return &rummy.StartGameResponse{}, err
}

func (s *RummyServer) SubscribeGame(req *rummy.SubscribeGameRequest, stream rummy.RummyService_SubscribeGameServer) error {
	glog.V(1).Infof("SubscribeGame: %v", req)
	s.gamesMu.Lock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return fmt.Errorf("no such game: %v", req.GameName)
	}

	eventsCh := make(chan *rummy.GameEvent, eventsBufferSize)
	g.Subscribe(eventsCh)
	s.gamesMu.Unlock()

	for e := range eventsCh {
		if err := stream.Send(e); err != nil {
			return err
		}
	}

	return nil
}

func (s *RummyServer) GetGameState(ctx context.Context, req *rummy.GetGameStateRequest) (*rummy.GameState, error) {
	glog.V(1).Infof("GetGameState: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	return g.GameState(), nil
}

func protoSlice(cards []deck.Card) []*deck.Card {
	result := make([]*deck.Card, len(cards))
	for i := range cards {
		result[i] = &cards[i]
	}
	return result
}

func valueSlice(cards []*deck.Card) []deck.Card {
	result := make([]deck.Card, len(cards))
	for i, c := range cards {
		result[i] = *c
	}
	return result
}

func (s *RummyServer) GetHandCards(ctx context.Context, req *rummy.GetHandCardsRequest) (*rummy.GetHandCardsResponse, error) {
	glog.V(1).Infof("GetHandCards: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	cards, err := g.PlayerHand(req.PlayerId)
	if err != nil {
		return nil, err
	}

	glog.Infof("Getting current hand for player %v in game %v: %s",
		req.PlayerId, req.GameName, rummy.NewHand(cards))
	return &rummy.GetHandCardsResponse{
		Cards: protoSlice(cards),
	}, nil
}

func (s *RummyServer) PickUpStock(ctx context.Context, req *rummy.PickUpStockRequest) (*rummy.PickUpStockResponse, error) {
	glog.V(1).Infof("PickUpStock: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	card, err := g.PickUpStock(req.PlayerId)
	return &rummy.PickUpStockResponse{
		Card: &card,
	}, err
}

func (s *RummyServer) PickUpDiscard(ctx context.Context, req *rummy.PickUpDiscardRequest) (*rummy.PickUpDiscardResponse, error) {
	glog.V(1).Infof("PickUpDiscard: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	cards, err := g.PickUpDiscard(req.PlayerId, int(req.NCards))
	return &rummy.PickUpDiscardResponse{
		Cards: protoSlice(cards),
	}, err
}

func (s *RummyServer) PlayCards(ctx context.Context, req *rummy.PlayCardsRequest) (*rummy.PlayCardsResponse, error) {
	glog.V(1).Infof("PlayCards: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	score, err := g.PlayCards(req.PlayerId, valueSlice(req.Cards))
	return &rummy.PlayCardsResponse{
		Score: int32(score),
	}, err
}

func (s *RummyServer) DiscardCard(ctx context.Context, req *rummy.DiscardCardRequest) (*rummy.DiscardCardResponse, error) {
	glog.V(1).Infof("DiscardCard: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	err := g.DiscardCard(req.PlayerId, *req.Card)
	return &rummy.DiscardCardResponse{}, err
}

func (s *RummyServer) CallRummy(ctx context.Context, req *rummy.CallRummyRequest) (*rummy.CallRummyResponse, error) {
	glog.V(1).Infof("CallRummy: %v", req)
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	var g *rummy.Game
	var ok bool
	if g, ok = s.games[req.GameName]; !ok {
		return nil, fmt.Errorf("no such game: %v", req.GameName)
	}

	err := g.CallRummy(req.PlayerId, valueSlice(req.Cards))
	return &rummy.CallRummyResponse{}, err
}
