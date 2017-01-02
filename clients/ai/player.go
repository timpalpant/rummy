package ai

import (
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"rummy"
	"rummy/clients/ai/strategy"
	"rummy/deck"
)

// Play the given game with the given strategy.
// AIs play by issuing requests to the game server, just like human players.
// The given Strategy will be invoked at the appropriate times to make
// play decisions.
func PlayGame(connStr, game string, playerId int32, strategy strategy.Strategy) error {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rummy.NewRummyServiceClient(conn)

	p := &computerPlayer{client, game, playerId, strategy}
	return p.Play()
}

// computerPlayer automatically initiates gameplay actions when it
// is their turn according to a certain strategy.
type computerPlayer struct {
	client   rummy.RummyServiceClient
	gameName string
	playerId int32
	strategy strategy.Strategy
}

func (cp *computerPlayer) Play() error {
	req := &rummy.SubscribeGameRequest{
		GameName: cp.gameName,
	}
	stream, err := cp.client.SubscribeGame(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if resp.Type == rummy.GameEvent_TURN_START && resp.PlayerId == cp.playerId {
			if err := cp.playTurn(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cp *computerPlayer) playTurn() error {
	req := &rummy.GetGameStateRequest{
		GameName: cp.gameName,
	}
	gs, err := cp.client.GetGameState(context.Background(), req)
	if err != nil {
		return err
	}

	discardPile := make([]deck.Card, len(gs.DiscardPile))
	for i, c := range gs.DiscardPile {
		discardPile[i] = *c
	}

	pickUpDiscard := cp.strategy.PickUpCards(discardPile)
	if pickUpDiscard == 0 {
		_, err := cp.client.PickUpStock(context.Background(), &rummy.PickUpStockRequest{
			PlayerId: cp.playerId,
		})
		if err != nil {
			return err
		}
	} else {
		_, err := cp.client.PickUpDiscard(context.Background(), &rummy.PickUpDiscardRequest{
			PlayerId: cp.playerId,
			NCards:   int32(pickUpDiscard),
		})
		if err != nil {
			return err
		}
	}

	resp, err := cp.client.GetHandCards(context.Background(), &rummy.GetHandCardsRequest{
		PlayerId: cp.playerId,
	})
	hand := make(rummy.Hand, len(resp.Cards))
	for _, c := range resp.Cards {
		hand[*c] = struct{}{}
	}

	playCards := cp.strategy.PlayCards(hand)
	if len(playCards) > 0 {
		cards := make([]*deck.Card, len(playCards))
		for i, c := range playCards {
			cards[i] = &c
		}

		_, err := cp.client.PlayCards(context.Background(), &rummy.PlayCardsRequest{
			PlayerId: cp.playerId,
			Cards:    cards,
		})
		if err != nil {
			return err
		}
	}

	discardCard := cp.strategy.Discard(hand)
	_, err = cp.client.DiscardCard(context.Background(), &rummy.DiscardCardRequest{
		PlayerId: cp.playerId,
		Card:     &discardCard,
	})
	return err
}
