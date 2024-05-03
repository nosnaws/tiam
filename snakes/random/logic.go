package main

import (
	"log"
	"math/rand"
	"time"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
)

func determineMove(state api.GameState) b.Move {
	board := b.BuildBoard(state)
	moves := board.GetMovesForSnake(b.MeId)
	if len(moves) < 1 {
		return b.Left
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]b.SnakeMove, len(moves))
	perm := r.Perm(len(moves))
	for i, randIndex := range perm {
		ret[i] = moves[randIndex]
	}

	return ret[0].Dir
}

func info() api.BattlesnakeInfoResponse {
	log.Println("INFO")
	return api.BattlesnakeInfoResponse{
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
func start(state api.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state api.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state api.GameState) api.BattlesnakeMoveResponse {
	log.Println("START TURN: ", state.Turn)

	move := determineMove(state)

	log.Println("RETURNING TURN: ", state.Turn, move)
	if move == b.Left {
		return api.BattlesnakeMoveResponse{
			Move: "left",
		}
	}
	if move == b.Right {
		return api.BattlesnakeMoveResponse{
			Move: "right",
		}
	}
	if move == b.Up {
		return api.BattlesnakeMoveResponse{
			Move: "up",
		}
	}
	return api.BattlesnakeMoveResponse{
		Move: "down",
	}
}
