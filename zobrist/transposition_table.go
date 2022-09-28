package zobrist

import (
	"math/rand"
	"time"

	"github.com/nosnaws/tiam/board"
)

type Key uint64
type ZobristTable [][]Key

type piece int

const (
	numPieces int = 10

	emptyPiece  piece = -1
	foodPiece         = 0
	hazardPiece       = 1

	meTripleStack = 2
	meDoubleStack = 3
	meHead        = 4
	meSegment     = 5

	otherTripleStack = 6
	otherDoubleStack = 7
	otherHead        = 8
	otherSegment     = 9
)

func InitializeZobristTable(height, width int) ZobristTable {
	rand.Seed(time.Now().Unix())
	boardLen := height * width
	t := make(ZobristTable, boardLen)

	for i := 0; i < boardLen; i++ {
		t[i] = make([]Key, numPieces)
		for j := 0; j < numPieces; j++ {
			t[i][j] = Key(rand.Uint64())
		}
	}

	return t
}

func GetZobristKey(zh ZobristTable, b *board.FastBoard) Key {
	key := Key(0)
	boardLen := len(b.List)

	for i := 0; i < boardLen; i++ {
		p := indexOf(b.List[i])

		if p != emptyPiece {
			key ^= zh[i][p]
		}
	}

	return key
}

func indexOf(t board.Tile) piece {
	id, ok := t.GetSnakeId()
	if ok {
		if id == board.MeId {
			if t.IsTripleStack() {
				return meTripleStack
			} else if t.IsDoubleStack() {
				return meDoubleStack
			} else if t.IsSnakeHead() {
				return meHead
			} else if t.IsSnakeBodyPart() {
				return meSegment
			}
		} else {
			if t.IsTripleStack() {
				return otherTripleStack
			} else if t.IsDoubleStack() {
				return otherDoubleStack
			} else if t.IsSnakeHead() {
				return otherHead
			} else if t.IsSnakeBodyPart() {
				return otherSegment
			}
		}
	} else if t.IsFood() {
		return foodPiece
	} else if t.IsHazard() {
		return hazardPiece
	}

	return emptyPiece
}
