package main

import (
	"log"
	"sort"

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
	isWrapped := state.Game.Ruleset.Name == "wrapped"
	board := g.BuildBoard(state)
	moves := board.GetMovesForSnake(g.MeId)
	if len(moves) < 1 {
		return g.Left
	}

	movesScore := []moveAndScore{}
	for _, move := range moves {
		indexInDir := g.IndexInDirection(
			move.Dir,
			board.Heads[g.MeId],
			uint16(state.Board.Width),
			uint16(state.Board.Height),
			isWrapped,
		)

		// Short curcuit if we find food
		if board.IsTileFood(indexInDir) {
			return move.Dir
		}

		ns := board.Clone()
		m := make(map[g.SnakeId]g.Move)
		m[g.MeId] = move.Dir
		ns.AdvanceBoard(m)

		s := g.Voronoi(&ns, g.MeId)

		movesScore = append(movesScore, moveAndScore{
			move:         move.Dir,
			voronoiScore: int(s.Score[g.MeId]),
			foodScore:    s.FoodDepth[g.MeId],
		})
	}

	sort.Slice(movesScore, func(i, j int) bool {
		return compareMoves(movesScore[i], movesScore[j])
	})

	return movesScore[0].move
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
