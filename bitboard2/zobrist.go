package bitboard

import (
	"math/rand"
)

type Key uint64
type ZobristTable [][]Key

type piece int

const (
	numPieces int = 4

	emptyPiece    piece = -1
	foodPiece           = 0
	hazardPiece         = 1
	foodOnHazard        = 2
	occupiedPiece       = 3
)

func (bb *BitBoard) InitializeZobristTable() ZobristTable {
	boardLen := bb.height * bb.width
	t := make(ZobristTable, boardLen)

	for i := 0; i < boardLen; i++ {
		t[i] = make([]Key, numPieces)
		for j := 0; j < numPieces; j++ {
			t[i][j] = Key(rand.Uint64())
		}
	}

	return t
}

func (bb *BitBoard) GetZobristKey(zh ZobristTable) Key {
	key := Key(0)
	boardLen := bb.width * bb.height

	for i := 0; i < boardLen; i++ {
		p := bb.indexOf(i)

		if p != emptyPiece {
			key ^= zh[i][p]
		}
	}

	return key
}

func (bb *BitBoard) indexOf(i int) piece {
	if bb.IsIndexOccupied(i) {
		return occupiedPiece
	}

	if bb.IsIndexFood(i) && bb.IsIndexHazard(i) {
		return foodOnHazard
	}

	if bb.IsIndexFood(i) {
		return foodPiece
	}

	if bb.IsIndexHazard(i) {
		return hazardPiece
	}

	return emptyPiece
}
