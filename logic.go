package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	//"github.com/newrelic/go-agent/v3/newrelic"
	"log"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nosnaws/tiam/brain"
	fastGame "github.com/nosnaws/tiam/game"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() fastGame.BattlesnakeInfoResponse {
	log.Println("INFO")
	return fastGame.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",           // TODO: Your Battlesnake username
		Color:      "#002080",    // TODO: Personalize
		Head:       "evil",       // TODO: Personalize
		Tail:       "fat-rattle", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state fastGame.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state fastGame.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state fastGame.GameState, txn *newrelic.Transaction) fastGame.BattlesnakeMoveResponse {
	log.Println("START TURN: ", state.Turn)
	gameBoard := fastGame.BuildBoard(state)

	move := brain.MCTS(&gameBoard, txn)

	log.Println("RETURNING TURN: ", state.Turn)
	if move.Dir == fastGame.Left {
		return fastGame.BattlesnakeMoveResponse{
			Move: "left",
		}
	}
	if move.Dir == fastGame.Right {
		return fastGame.BattlesnakeMoveResponse{
			Move: "right",
		}
	}
	if move.Dir == fastGame.Up {
		return fastGame.BattlesnakeMoveResponse{
			Move: "up",
		}
	}
	return fastGame.BattlesnakeMoveResponse{
		Move: "down",
	}
}
