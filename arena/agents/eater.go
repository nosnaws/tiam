package agents

import (
	"math"
	"sort"

	bitboard "github.com/nosnaws/tiam/bitboard2"
)

type scoredMove struct {
	score float64
	move  bitboard.SnakeMoveSet
}

func GetEaterMove(bb *bitboard.BitBoard, snakeId string) bitboard.SnakeMoveSet {
	moves := bb.SplitMoves(bb.GetMoves(snakeId))

	scoredMoves := []scoredMove{}

	for _, m := range moves {
		ns := bb.Clone()
		ns.AdvanceTurn([]bitboard.SnakeMoveSet{m})

		scoredMoves = append(scoredMoves, scoredMove{
			move:  m,
			score: scoreMove(ns, snakeId),
		})
	}

	sort.Slice(scoredMoves, func(i, j int) bool {
		return scoredMoves[i].score > scoredMoves[j].score
	})

	return scoredMoves[0].move
}

func scoreMove(bb *bitboard.BitBoard, snakeId string) float64 {
	if !bb.IsSnakeAlive(snakeId) {
		return -math.MaxFloat64
	}

	food := bb.FoodScore(snakeId)
	snake := bb.GetSnake(snakeId)
	length := float64(snake.Length) * 0.5

	return food + length
}
