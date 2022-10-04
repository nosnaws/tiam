package main

import (
	"context"
	"log"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
	"github.com/nosnaws/tiam/mmm"
	min "github.com/nosnaws/tiam/mmm"
)

var gameCache map[string]*mmm.Cache

type moveAndScore struct {
	move         b.Move
	voronoiScore int
	foodScore    int
}

func initialize() {
	gameCache = make(map[string]*min.Cache)
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

func determineMove(ctx context.Context, state api.GameState) b.Move {
	board := b.BuildBoard(state)
	cache := gameCache[state.Game.ID]
	cache.SetCurTurn(state.Turn)
	//move := b.Minmax(&board, g.MeId, 6)
	//move := b.IdfsMinmax(&board)
	//move := b.BRS(&board, g.Left, 8, math.Inf(-1), math.Inf(1), true)
	//move := b.IDBRS(ctx, &board)
	//move, _ := min.MultiMinmax(&board, 12)
	move := min.MultiMinmaxID(&board, cache)

	return move
}

func info() api.BattlesnakeInfoResponse {
	log.Println("INFO")
	return api.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",
		Color:      "#a57aa8",
		Head:       "rocket-helmet",
		Tail:       "rocket",
	}
}

func start(state api.GameState) {
	board := b.BuildBoard(state)
	gameCache[state.Game.ID] = mmm.CreateCache(&board, 0)
	log.Printf("%s START\n", state.Game.ID)
}

func end(state api.GameState) {
	delete(gameCache, state.Game.ID)
	log.Printf("%s END\n\n", state.Game.ID)
}

func move(ctx context.Context, state api.GameState) api.BattlesnakeMoveResponse {
	log.Println("START TURN: ", state.Turn)

	move := determineMove(ctx, state)

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
