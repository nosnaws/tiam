package bitboard

import (
	"math/rand"

	"github.com/nosnaws/tiam/moveset"
)

func (bb *BitBoard) randomMove(snakeId string, rand *rand.Rand) SnakeMoveSet {
	if len(bb.GetSnake(snakeId).body) == 0 {
		bb.Print()
		panic("why????")
	}
	moves := bb.GetMoves(snakeId)

	//if len(moves) == 0 {
	//return SnakeMove{Id: snakeId, Dir: Left}
	//}

	return SnakeMoveSet{
		Id:  snakeId,
		Set: moveset.GetRandomMove(moves.Set, rand),
	}
}

func (bb *BitBoard) RandomPlayout(length int, rand *rand.Rand) {
	maxTurn := length + bb.turn

	for !bb.IsGameOver() && bb.turn < maxTurn {
		moves := []SnakeMoveSet{}
		for id, snake := range bb.Snakes {
			if snake.IsAlive() {
				moves = append(moves, bb.randomMove(id, rand))
			}
		}
		bb.AdvanceTurn(moves)
	}
}
