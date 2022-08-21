package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nosnaws/tiam/brain"
	g "github.com/nosnaws/tiam/game"
)

func determineMove(state g.GameState) g.Move {
	board := g.BuildBoard(state)
	moves := board.GetMovesForSnake(g.MeId)
	if len(moves) < 1 {
		return g.Left
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]g.SnakeMove, len(moves))
	perm := r.Perm(len(moves))
	for i, randIndex := range perm {
		ret[i] = moves[randIndex]
	}

	return ret[0].Dir
}

func info() g.BattlesnakeInfoResponse {
	log.Println("INFO")
	return g.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",   // TODO: Your Battlesnake username
		Color:      "#000000",   // TODO: Personalize
		Head:       "bendr",     // TODO: Personalize
		Tail:       "block-bum", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state g.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state g.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state g.GameState, config *brain.MCTSConfig, txn *newrelic.Transaction) g.BattlesnakeMoveResponse {
	log.Println("START TURN: ", state.Turn)

	move := determineMove(state)

	log.Println("RETURNING TURN: ", state.Turn, move)
	if move == g.Left {
		return g.BattlesnakeMoveResponse{
			Move: "left",
		}
	}
	if move == g.Right {
		return g.BattlesnakeMoveResponse{
			Move: "right",
		}
	}
	if move == g.Up {
		return g.BattlesnakeMoveResponse{
			Move: "up",
		}
	}
	return g.BattlesnakeMoveResponse{
		Move: "down",
	}
}
