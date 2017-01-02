// cli is a command-line interface to the rummy server.
// Connect with `./cli -server host:port`.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"rummy"
	"rummy/deck"
)

var stdin = bufio.NewReader(os.Stdin)

func prompt(msg string) string {
	fmt.Print(msg)
	result, err := stdin.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.TrimRight(result, "\n")
}

func createGame(client rummy.RummyServiceClient) (string, int32, error) {
	gameName := prompt("Enter game name: ")
	_, err := client.CreateGame(context.Background(), &rummy.CreateGameRequest{
		GameName: gameName,
	})
	if err != nil {
		return "", 0, err
	}

	playerName := prompt("Enter player name: ")
	resp, err := client.JoinGame(context.Background(), &rummy.JoinGameRequest{
		GameName:   gameName,
		PlayerName: playerName,
	})
	if err != nil {
		return "", 0, err
	}
	playerId := resp.PlayerId

	var start string
	for start != "y" {
		resp, err := client.GetGameState(context.Background(), &rummy.GetGameStateRequest{
			GameName: gameName,
		})
		if err != nil {
			return "", 0, err
		}

		fmt.Println("Current players in game:")
		for _, player := range resp.Players {
			fmt.Printf("\t%v: %v\n", player.Id, player.Name)
		}

		start = prompt("Start game? (y/n): ")
	}

	_, err = client.StartGame(context.Background(), &rummy.StartGameRequest{
		GameName: gameName,
	})

	return gameName, playerId, err
}

func joinGame(client rummy.RummyServiceClient) (string, int32, error) {
	gameName := prompt("Enter game name: ")
	playerName := prompt("Enter player name: ")
	strategyName := prompt("Strategy? (leave blank to play interactively): ")

	resp, err := client.JoinGame(context.Background(), &rummy.JoinGameRequest{
		GameName:   gameName,
		PlayerName: playerName,
		Strategy:   strategyName,
	})
	if err != nil {
		return "", 0, err
	}

	fmt.Println("Waiting for game to start")
	nPlayers := 0
	for {
		resp, err := client.GetGameState(context.Background(), &rummy.GetGameStateRequest{
			GameName: gameName,
		})
		if err != nil {
			return "", 0, err
		}

		if resp.CurrentPlayerTurn != -1 {
			fmt.Println("Game has started")
			break
		}

		if len(resp.Players) != nPlayers {
			fmt.Println("Current players in game:")
			for _, player := range resp.Players {
				fmt.Printf("\t%v: %v\n", player.Id, player.Name)
			}
			nPlayers = len(resp.Players)
		}

		time.Sleep(2 * time.Second)
	}

	return gameName, resp.PlayerId, nil
}

func playGame(client rummy.RummyServiceClient, gameName string, playerId int32) {
	resp, err := client.GetGameState(context.Background(), &rummy.GetGameStateRequest{
		GameName: gameName,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	playerNames := make([]string, len(resp.Players))
	for _, p := range resp.Players {
		playerNames[p.Id] = p.Name
	}

	req := &rummy.SubscribeGameRequest{
		GameName: gameName,
	}
	stream, err := client.SubscribeGame(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.CurrentPlayerTurn == playerId {
		// We're up first, play our first turn.
		if err := playTurn(client, gameName, playerId); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("%v's turn\n", playerNames[resp.CurrentPlayerTurn])
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		printGameEvent(resp, playerNames)
		if resp.Type == rummy.GameEvent_TURN_START && resp.PlayerId == playerId {
			if err := playTurn(client, gameName, playerId); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	printEndGame(client, gameName)
}

func printGameEvent(e *rummy.GameEvent, playerNames []string) {
	s := playerNames[e.PlayerId]

	switch e.Type {
	case rummy.GameEvent_TURN_START:
		s = s + "'s turn."
	case rummy.GameEvent_PICK_UP_STOCK:
		s = s + " picked up a card from the stock."
	case rummy.GameEvent_PICK_UP_DISCARD:
		s = s + " picked up from the discard pile"
	case rummy.GameEvent_PLAY_CARDS:
		s = s + " played cards"
	case rummy.GameEvent_DISCARD:
		s = s + " discarded"
	}

	if len(e.Cards) > 0 {
		s = s + fmt.Sprintf(": %v", ppCards(e.Cards, false))
	}

	if e.Score != 0 {
		s = s + fmt.Sprintf(" for %d points", e.Score)
	}

	fmt.Println(s)
}

func printCurrentHand(client rummy.RummyServiceClient, gameName string, playerId int32) error {
	resp, err := client.GetHandCards(context.Background(), &rummy.GetHandCardsRequest{
		GameName: gameName,
		PlayerId: playerId,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Current hand: %v\n", ppCards(resp.Cards, true))
	return nil
}

func printCurrentDiscardPile(client rummy.RummyServiceClient, gameName string) error {
	gs, err := client.GetGameState(context.Background(), &rummy.GetGameStateRequest{
		GameName: gameName,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Current discard pile: %v\n", ppCards(gs.DiscardPile, false))
	return nil
}

func playTurn(client rummy.RummyServiceClient, gameName string, playerId int32) error {
	fmt.Println("\nYour turn!")
	if err := printCurrentHand(client, gameName, playerId); err != nil {
		return err
	}
	if err := printCurrentDiscardPile(client, gameName); err != nil {
		return err
	}

	pickUpCards(client, gameName, playerId)

	for {
		playCards(client, gameName, playerId)
		if err := discard(client, gameName, playerId); err == nil {
			break
		} else {
			fmt.Printf("Error discarding: %v\n", err)
		}
	}

	if err := printCurrentHand(client, gameName, playerId); err != nil {
		return err
	}
	if err := printCurrentDiscardPile(client, gameName); err != nil {
		return err
	}

	return nil
}

func pickUpCards(client rummy.RummyServiceClient, gameName string, playerId int32) {
	for {
		fmt.Println("\nWhat would you like to do?")
		fmt.Println("\t1) Pick up a card from the stock")
		fmt.Println("\t2) Pick up card(s) from the discard pile")
		choice := prompt("Selection: ")
		switch choice {
		case "1":
			if err := pickUpStock(client, gameName, playerId); err != nil {
				fmt.Printf("Error picking up from stock: %v\n", err)
			} else {
				return
			}
		case "2":
			if err := pickUpDiscard(client, gameName, playerId); err != nil {
				fmt.Printf("Error picking up from discard: %v\n", err)
			} else {
				return
			}
		default:
			fmt.Println("Invalid selection, please choose 1 or 2")
		}
	}
}

func pickUpStock(client rummy.RummyServiceClient, gameName string, playerId int32) error {
	resp, err := client.PickUpStock(context.Background(), &rummy.PickUpStockRequest{
		GameName: gameName,
		PlayerId: playerId,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Picked up: %v\n", deck.CardString(*resp.Card))
	return nil
}

func pickUpDiscard(client rummy.RummyServiceClient, gameName string, playerId int32) error {
	n := 0
	var err error
	for n == 0 {
		nStr := prompt("How many cards would you like to pick up?: ")
		n, err = strconv.Atoi(nStr)
		if err != nil {
			fmt.Printf("Invalid number '%v': %v", nStr, err)
		}
	}

	resp, err := client.PickUpDiscard(context.Background(), &rummy.PickUpDiscardRequest{
		GameName: gameName,
		PlayerId: playerId,
		NCards:   int32(n),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Picked up: %v\n", ppCards(resp.Cards, false))
	return nil
}

func playCards(client rummy.RummyServiceClient, gameName string, playerId int32) {
	for {
		resp, err := client.GetHandCards(context.Background(), &rummy.GetHandCardsRequest{
			GameName: gameName,
			PlayerId: playerId,
		})
		if err != nil {
			fmt.Printf("Error getting current hand: %v\n", err)
			continue
		}
		sort.Sort(bySuitAndRank(resp.Cards))
		numbered := make([]string, len(resp.Cards))
		for i, c := range resp.Cards {
			numbered[i] = fmt.Sprintf("%d:%v", i, deck.CardString(*c))
		}
		fmt.Printf("Current hand: %v\n", strings.Join(numbered, " "))

		cardsToPlay := prompt(
			"Select cards to play as a meld or a rummy (e.g. 1,4,5). " +
				"Leave empty to continue: ")
		if cardsToPlay == "" {
			break
		}

		cards, err := parseCardSelection(resp.Cards, cardsToPlay)
		if err != nil {
			fmt.Printf("Cannot select '%v': %v\n", cardsToPlay, err)
			continue
		}

		playResp, err := client.PlayCards(context.Background(), &rummy.PlayCardsRequest{
			GameName: gameName,
			PlayerId: playerId,
			Cards:    cards,
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Played %v for %v points\n", ppCards(cards, true), playResp.Score)
		}
	}
}

func parseCardSelection(hand []*deck.Card, selectionStr string) ([]*deck.Card, error) {
	cardStrs := strings.Split(selectionStr, ",")
	result := make([]*deck.Card, len(cardStrs))
	for i, cardStr := range cardStrs {
		cardIdx, err := strconv.Atoi(strings.TrimSpace(cardStr))
		if err != nil {
			return nil, fmt.Errorf("error parsing %v: %v", cardStr, err)
		} else if cardIdx < 0 || cardIdx >= len(hand) {
			return nil, fmt.Errorf("invalid card selection: %v", cardStr)
		}
		result[i] = hand[cardIdx]
	}

	return result, nil
}

func discard(client rummy.RummyServiceClient, gameName string, playerId int32) error {
	resp, err := client.GetHandCards(context.Background(), &rummy.GetHandCardsRequest{
		GameName: gameName,
		PlayerId: playerId,
	})
	if err != nil {
		return fmt.Errorf("Error getting current hand: %v", err)
	}
	numbered := make([]string, len(resp.Cards))
	for i, c := range resp.Cards {
		numbered[i] = fmt.Sprintf("%d:%v", i, deck.CardString(*c))
	}
	fmt.Printf("Current hand: %v\n", strings.Join(numbered, " "))

	cardToDiscard := prompt("Select card to discard: ")
	cards, err := parseCardSelection(resp.Cards, cardToDiscard)
	if err != nil || len(cards) != 1 {
		return fmt.Errorf("Cannot select '%v': %v", cardToDiscard, err)
	}

	_, err = client.DiscardCard(context.Background(), &rummy.DiscardCardRequest{
		GameName: gameName,
		PlayerId: playerId,
		Card:     cards[0],
	})
	return err
}

// Sort cards by suit and then rank.
type bySuitAndRank []*deck.Card

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

func ppCards(cards []*deck.Card, sorted bool) string {
	if sorted {
		sort.Sort(bySuitAndRank(cards))
	}

	result := make([]string, len(cards))
	for i, c := range cards {
		result[i] = deck.CardString(*c)
	}

	return fmt.Sprintf("%v", result)
}

func printEndGame(client rummy.RummyServiceClient, gameName string) {
	resp, err := client.GetGameState(context.Background(), &rummy.GetGameStateRequest{
		GameName: gameName,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.GameOver {
		fmt.Println("ERROR: Game ended prematurely, please re-join")
		return
	}

	sort.Sort(byScore(resp.Players))
	winner := resp.Players[0]
	fmt.Printf("Game over: %v wins!\n", winner.Name)
	fmt.Println("Final scores:")
	for _, player := range resp.Players {
		fmt.Printf("\t%v: %v\n", player.Name, player.CurrentScore)
	}
}

// byScore sorts players by descending score.
// The winner at the end of the game will be byScore[0].
type byScore []*rummy.PlayerState

func (b byScore) Len() int {
	return len(b)
}

func (b byScore) Less(i, j int) bool {
	return b[i].CurrentScore > b[j].CurrentScore
}

func (b byScore) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func printMainMenu() {
	fmt.Println("\nMain menu:")
	fmt.Println("\t1) Create a new game")
	fmt.Println("\t2) Join a game")
	fmt.Println("\t3) Quit")
}

func main() {
	connStr := flag.String("server", "", "Game server to connect to")
	flag.Parse()

	if *connStr == "" {
		fmt.Println("You must specify a server to connect to with -server")
		os.Exit(1)
	}

	conn, err := grpc.Dial(*connStr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	client := rummy.NewRummyServiceClient(conn)

	fmt.Println("Welcome to Rummy!")
	for {
		printMainMenu()
		selection := prompt("Please make a selection: ")

		switch selection {
		case "1":
			gameName, playerId, err := createGame(client)
			if err != nil {
				fmt.Println(err)
				continue
			}
			playGame(client, gameName, playerId)
		case "2":
			gameName, playerId, err := joinGame(client)
			if err != nil {
				fmt.Println(err)
				continue
			}
			playGame(client, gameName, playerId)
		case "3":
			return
		}
	}
}
