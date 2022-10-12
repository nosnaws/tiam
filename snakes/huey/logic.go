package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
	"github.com/nosnaws/tiam/mmm"
	min "github.com/nosnaws/tiam/mmm"
)

//var gameCache map[string]*mmm.Cache

var CONFIG mmm.GMOConfig

type moveAndScore struct {
	move         b.Move
	voronoiScore int
	foodScore    int
}

func initialize() {
	//gameCache = make(map[string]*min.Cache)
	initializeGMO()
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

func initializeGMO() {
	CONFIG = mmm.GMOConfig{
		FoodWeightA:    getEnv("FOOD_A", 1),
		FoodWeightB:    getEnv("FOOD_B", 1),
		VoronoiWeightA: getEnv("VORONOI_A", 1),
		VoronoiWeightB: getEnv("VORONOI_B", 1),
		LengthWeight:   getEnv("LENGTH", 1),
	}
}

func determineMove(state api.GameState) b.Move {
	board := b.BuildBoard(state)
	//cache := gameCache[state.Game.ID]
	//cache.SetCurTurn(state.Turn)
	//move := b.Minmax(&board, g.MeId, 6)
	//move := b.IdfsMinmax(&board)
	//move := b.BRS(&board, g.Left, 8, math.Inf(-1), math.Inf(1), true)
	//move := b.IDBRS(ctx, &board)
	//move, _ := min.MultiMinmax(&board, 12)
	evoStrat := mmm.GMOConfig{
		FoodWeightA:    4.166351,
		FoodWeightB:    3.785205,
		VoronoiWeightA: 1.900364,
		VoronoiWeightB: 1.150082,
		LengthWeight:   2.334582,
	}
	//strat := mmm.CreateGMOStrategy(CONFIG)
	move, score := min.MultiMMLimited(&board, 2, mmm.CreateGMOStrategy(evoStrat))
	fmt.Println("SCORE", score)

	return move
}

func info() api.BattlesnakeInfoResponse {
	log.Println("INFO")
	return api.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nosnaws",
		Color:      "#349eeb",
		Head:       "lantern-fish",
		Tail:       "mlh-gene",
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
