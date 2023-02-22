package main

import (
	"context"
	"log"
	"os"
	"strconv"

	api "github.com/nosnaws/tiam/battlesnake"
	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/paranoid"
)

func getEnv(name string, defaultVal float64) float64 {
	if val, ok := os.LookupEnv(name); ok {
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatalln("Unable to parse env", name, val)
		}
		return num
	}

	return defaultVal
}

type moveAndScore struct {
	dir   bitboard.Dir
	score float64
}

func determineMove(state api.GameState) bitboard.Dir {
	board := bitboard.CreateBitBoard(state)

	move := paranoid.GetMoveID(board, state.You.ID)
	return move
}

func info() api.BattlesnakeInfoResponse {
	log.Println("INFO")
	return api.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",
		Color:      "#3d5a80",
		Head:       "nr-rocket",
		Tail:       "comet",
	}
}

func start(state api.GameState) {
	//board := b.BuildBoard(state)
	//gameCache[state.Game.ID] = mmm.CreateCache(&board, 0)
	log.Printf("%s START\n", state.Game.ID)
}

func end(state api.GameState) {
	//delete(gameCache, state.Game.ID)
	log.Printf("%s END\n\n", state.Game.ID)
}

func move(ctx context.Context, state api.GameState) api.BattlesnakeMoveResponse {
	log.Println("START TURN: ", state.Turn)

	move := determineMove(state)

	log.Println("RETURNING TURN: ", state.Turn, move)
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
