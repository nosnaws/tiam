package bitboard

import (
	"math/rand"

	"github.com/nosnaws/tiam/moveset"
)

func (bb *BitBoard) RandomMove(snakeId string, rand *rand.Rand) SnakeMoveSet {
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
		Set: moveset.GetRandomMove(moves.Set),
	}
}

func (bb *BitBoard) RandomPlayout(length int, rand *rand.Rand) {
	maxTurn := length + bb.Turn

	for !bb.IsGameOver() && bb.Turn < maxTurn {
		moves := []SnakeMoveSet{}
		for id, snake := range bb.Snakes {
			if snake.IsAlive() {
				moves = append(moves, bb.RandomMove(id, rand))
			}
		}
		bb.AdvanceTurn(moves)
	}
}

func (bb *BitBoard) RandomPlayoutMonte(length int, rand *rand.Rand) Dir {
	maxTurn := length + bb.Turn
	startingTurn := bb.Turn
	firstMove := Dir("")

	for !bb.IsGameOver() && bb.Turn < maxTurn {
		moves := []SnakeMoveSet{}
		for id, snake := range bb.Snakes {
			if snake.IsAlive() {
				move := bb.RandomMove(id, rand)
				moves = append(moves, move)
				if id == bb.meId && bb.Turn == startingTurn {
					firstMove = MoveSetToDir(move.Set)
				}
			}
		}
		bb.AdvanceTurn(moves)
	}

	return firstMove
}
