package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	//"github.com/newrelic/go-agent/v3/newrelic"
	"fmt"
	"log"

	"github.com/newrelic/go-agent/v3/newrelic"

	api "github.com/nosnaws/tiam/battlesnake"
	bitboard "github.com/nosnaws/tiam/bitboard2"
	mctsv3 "github.com/nosnaws/tiam/mcts_v3"
	"github.com/nosnaws/tiam/mmm"
)

func info() api.BattlesnakeInfoResponse {
	log.Println("INFO")
	return api.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",
		Color:      "#002080",
		Head:       "alligator",
		Tail:       "cosmic-horror",
	}
}

var gameCache map[string]*mmm.Cache

func initialize() {
	gameCache = make(map[string]*mmm.Cache)
}

func start(gc *mctsv3.GameController, state api.GameState) {
	log.Printf("%s START\n", state.Game.ID)
	gc.StartGame(state)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(gc *mctsv3.GameController, state api.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
	gc.EndGame(state)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(gc *mctsv3.GameController, gameState api.GameState, txn *newrelic.Transaction) api.BattlesnakeMoveResponse {
	log.Println("START TURN: ", gameState.Turn)
	//gameBoard.Print()

	mctsConfig := mctsv3.MCTSConfig{
		MinimumSims:    10,
		SimLimit:       15,
		TreeDepthLimit: 50,
	}
	move := gc.GetNextMove(gameState, mctsConfig)

	fmt.Println("RETURNING TURN: ", gameState.Turn)
	if move == bitboard.Left {
		return api.BattlesnakeMoveResponse{
			Move: "left",
		}
	}
	if move == bitboard.Right {
		return api.BattlesnakeMoveResponse{
			Move: "right",
		}
	}
	if move == bitboard.Up {
		return api.BattlesnakeMoveResponse{
			Move: "up",
		}
	}
	return api.BattlesnakeMoveResponse{
		Move: "down",
	}
}
