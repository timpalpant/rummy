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

AI
--

Package `clients/ai` provides a simplified interface for implementing strategies that
can play the game. To implement a new strategy, implement the `Strategy` interface and
then add it to the registry in `clients/ai/strategy/registry.go`. A rudimentary greedy
strategy is implemented in `clients/ai/strategy/greedy.go`.

Two or more strategies can be played against each other a large number of times using
the driver in `clients/ai/battle`.

$ ./battle -strategies nop,greedy -num_games 1000 -seed 123

Building
--------

Regenerate the protos:

$ make proto

Build the server:

$ make server

Build the CLI:

$ make clients
