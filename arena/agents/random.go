package agents

import (
	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/moveset"
)

func GetRandomMove(bb *bitboard.BitBoard, snakeId string) bitboard.SnakeMoveSet {
	moves := bb.GetMoves(snakeId)

	return bitboard.SnakeMoveSet{
		Id:  snakeId,
		Set: moveset.GetRandomMove(moves.Set),
	}
}
