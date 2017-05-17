Rummy 500
=========

This repository contains a [gRPC](http://www.grpc.io/) server that manages games of [Rummy 500](https://en.wikipedia.org/wiki/500_rum),
and CLI and AI clients to connect to it. It was implemented over Christmas 2016 to prolong
a family tournament with remote play, and as a way to learn about gRPC.

Server
------

The message protocol for the game server is defined in `service.proto`. The service provides
APIs to create and join games, and to observe and play in games you have joined. Multiple
games can be played simultaneously. A primitive form of authentication is provided by allowing
players to join with a provided secret that must be used for all subsequent gameplay.

During a game, clients subscribe to game events to observe the play of others and to know
when it is their turn. Game play events are pushed to subscribed clients using gRPC streaming.

Via [gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway), the server supports
the same API over REST/JSON:

Game management
- POST /v1/create/{game_name}
- POST /v1/join/{game_name}/{player_name}
- POST /v1/start/{game_name}

Game observation
- GET /v1/subscribe/{game_name}
- GET /v1/state/{game_name}
- POST /v1/hand

Game play
- POST /v1/pick_up_stock
- POST /v1/pick_up_discard
- POST /v1/play_cards
- POST /v1/discard
- POST /v1/call_rummy

This enables future development of a web interface.

State machine
-------------

Gameplay proceeds according to a state machine (see `GameState` in `game.proto`).
The current player begins in state `TURN_START` and must issue either a `PickUpStockRequest`
or a `PickUpDiscardRequest` to begin their turn. Once they have successfully picked up
cards with a valid request, the turn state transitions to `PICKED_UP_CARDS`. The player
may then play cards for points (if possible) by issuing a `PlayCardsRequest`, which
transitions their state to `PLAYED_CARDS`. Finally, the player must discard a card with a
`DiscardCardRequest` to end their turn.

The `Game` struct maintains the game invariants imposed by the rules of Rummy 500,
such as enforcing that after a player has picked up from the discard they must play
the bottom card for points before ending their turn.

CLI
---

`clients/cli` provides a command-line client that can be used to connect and play games
interactively against other human or AI players.

```
$ ./cli -server :8081
Welcome to Rummy!

Main menu:
	1) Create a new game
	2) Join a game
	3) Quit
Please make a selection: 1
Enter game name: TestGame
Enter player name: Tim
Add CP? (y/n): y
Enter strategy name: greedy
Add CP? (y/n): n
Current players in game:
	0: Tim
	1: CP0-greedy
Start game? (y/n): y

Your turn!
Current hand: [3♥ 5♥ 7♦ 9♦ 7♣ 9♠ J♠]
Current discard pile: [5♣ 7♠]
All played melds:
	[10♥ 10♣ 10♠]
Current player status:
	Tim: 7 cards, 0 points
	CP0-greedy: 4 cards, 30 points

What would you like to do?
	1) Pick up a card from the stock
	2) Pick up card(s) from the discard pile
Selection: 2
How many cards would you like to pick up?: 3
Error picking up from discard: rpc error: code = Unknown desc = can't pick up 3 > 2 cards in discard pile

What would you like to do?
	1) Pick up a card from the stock
	2) Pick up card(s) from the discard pile
Selection: 1
Picked up: A♦
Current hand: 0:3♥ 1:5♥ 2:A♦ 3:7♦ 4:9♦ 5:7♣ 6:9♠ 7:J♠
Select cards to play as a meld or a rummy (e.g. 1,4,5). Leave empty to continue:
Current hand: 0:3♥ 1:5♥ 2:A♦ 3:7♦ 4:9♦ 5:7♣ 6:9♠ 7:J♠
Select card to discard: 7
CP0-greedy's turn.
CP0-greedy picked up a card from the stock.
CP0-greedy discarded: [6♦]
```

AI
--

Package `clients/ai` provides a simplified interface for implementing strategies that
can play the game. To implement a new strategy, implement the `Strategy` interface and
then add it to the registry in `clients/ai/strategy/registry.go`. A rudimentary greedy
strategy is implemented in `clients/ai/strategy/greedy.go`.

Two or more strategies can be played against each other a large number of times using
the driver in `clients/ai/battle`.

> ./battle -strategies nop,greedy -num_games 1000 -seed 123

Building
--------

Regenerate the protos:
> make proto

Build the server:

> make server

Build the CLI:

> make clients
