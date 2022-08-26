package main

import (
	"log"
	"math"

	b "github.com/nosnaws/tiam/brain"
	g "github.com/nosnaws/tiam/game"
)

type moveAndScore struct {
	move         g.Move
	voronoiScore int
	foodScore    int
}

func compareMoves(a, b moveAndScore) bool {
	aFood := a.foodScore
	bFood := b.foodScore
	if aFood < 0 {
		aFood = 0
	}
	if bFood < 0 {
		bFood = 0
	}

	aScore := a.voronoiScore / (aFood + 1)
	bScore := b.voronoiScore / (bFood + 1)

	return aScore > bScore
}

func determineMove(state g.GameState) g.Move {
	board := g.BuildBoard(state)
	move := b.Minmax(&board, g.Move(""), 6, -math.MaxFloat64, math.MaxFloat64, true)

	return move.Move
}

func info() g.BattlesnakeInfoResponse {
	log.Println("INFO")
	return g.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",
		Color:      "#32a852",
		Head:       "iguana",
		Tail:       "iguana",
	}
}

func start(state g.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

func end(state g.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

func move(state g.GameState) g.BattlesnakeMoveResponse {
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
