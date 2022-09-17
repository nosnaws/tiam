package main

import (
	"context"
	"log"

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

func determineMove(ctx context.Context, state g.GameState) g.Move {
	board := g.BuildBoard(state)
	//move := b.Minmax(&board, g.MeId, 6)
	//move := b.IdfsMinmax(&board)
	//move := b.BRS(&board, g.Left, 8, math.Inf(-1), math.Inf(1), true)
	move := b.IDBRS(ctx, &board)

	return move
}

func info() g.BattlesnakeInfoResponse {
	log.Println("INFO")
	return g.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",
		Color:      "#38a852",
		Head:       "cosmic-horror",
		Tail:       "cosmic-horror",
	}
}

func start(state g.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

func end(state g.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

func move(ctx context.Context, state g.GameState) g.BattlesnakeMoveResponse {
	log.Println("START TURN: ", state.Turn)

	move := determineMove(ctx, state)

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
