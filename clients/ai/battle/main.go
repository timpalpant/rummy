package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/golang/glog"

	"rummy"
	"rummy/clients/ai"
	"rummy/clients/ai/strategy"
)

func main() {
	strategies := flag.String("strategies", "", "Strategies to simulate")
	numGames := flag.Int("num_games", 1000, "Number of games to simulate")
	seed := flag.Int64("seed", 1, "Seed for random shuffling")
	flag.Parse()

	rand.Seed(*seed)
	stratNames := strings.Split(*strategies, ",")
	if len(stratNames) < 2 {
		fmt.Println("You must specify strategies to simulate with -strategies")
		os.Exit(1)
	}

	results := make(map[string]int, len(stratNames))
	for _, stratName := range stratNames {
		results[stratName] = 0
	}

	glog.Infof("Simulating %v games", *numGames)
	for i := 0; i < *numGames; i++ {
		g := rummy.NewGame()
		id2StratName := make(map[int32]string, len(stratNames))
		for j, s := range stratNames {
			id, err := g.AddPlayer(fmt.Sprintf("CP%d", j))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			id2StratName[id] = s

			strat, err := strategy.ForName(s)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go func(s string) {
				if err := ai.PlayGame(g, id, strat); err != nil {
					glog.Fatal("Strategy %v errored: %v", s, err)
				}
			}(s)
		}

		// Start the game.
		if err := g.Deal(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Wait for game to finish.
		events := make(chan *rummy.GameEvent, 1)
		g.Subscribe(events)
		for ok := true; ok; _, ok = <-events {
		}

		// Collect game results for statistics.
		gs := g.GameState()
		var winner *rummy.PlayerState
		var winnerScore int32
		for _, p := range gs.Players {
			if winner == nil || p.CurrentScore > winnerScore {
				winner = p
			}
		}

		stratName := id2StratName[winner.Id]
		results[stratName]++
	}

	fmt.Println("\nResults:")
	for stratName, gamesWon := range results {
		fmt.Printf("%v: %v games won (%.2f %%)\n", stratName, gamesWon,
			float64(gamesWon)/float64(*numGames))
	}
}
