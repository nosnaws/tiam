package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	//"github.com/newrelic/go-agent/v3/newrelic"
	"log"

	"github.com/newrelic/go-agent/v3/newrelic"

	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/nosnaws/tiam/board"
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

func start(state api.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

func end(state api.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

func move(gameState api.GameState, txn *newrelic.Transaction) api.BattlesnakeMoveResponse {
	log.Println("START TURN: ", gameState.Turn)
	gameBoard := board.BuildBoard(gameState)

	move := mmm.MultiMinmaxID(&gameBoard)

	log.Println("RETURNING TURN: ", gameState.Turn)
	if move == board.Left {
		return api.BattlesnakeMoveResponse{
			Move: "left",
		}
	}
	if move == board.Right {
		return api.BattlesnakeMoveResponse{
			Move: "right",
		}
	}
	if move == board.Up {
		return api.BattlesnakeMoveResponse{
			Move: "up",
		}
	}
	return api.BattlesnakeMoveResponse{
		Move: "down",
	}
}
