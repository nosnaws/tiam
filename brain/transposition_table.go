package brain

import (
	"math/rand"
	"time"

	g "github.com/nosnaws/tiam/game"
)

type BoardHash uint64
type TranspositionTable [][]BoardHash

type Piece int

const (
	numPieces int = 10

	emptyPiece  Piece = -1
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

func InitializeTranspositionTable(height, width int) TranspositionTable {
	rand.Seed(time.Now().Unix())
	boardLen := height * width
	t := make(TranspositionTable, boardLen)

	for i := 0; i < boardLen; i++ {
		t[i] = make([]BoardHash, numPieces)
		for j := 0; j < numPieces; j++ {
			t[i][j] = BoardHash(rand.Uint64())
		}
	}

	return t
}

func HashBoard(tt TranspositionTable, b g.FastBoard) BoardHash {
	hash := BoardHash(0)
	boardLen := len(b.List)

	for i := 0; i < boardLen; i++ {
		p := indexOf(b.List[i])

		if p != emptyPiece {
			hash ^= tt[i][p]
		}
	}

	return hash
}

func indexOf(t g.Tile) Piece {
	id, ok := t.GetSnakeId()
	if ok {
		if id == g.MeId {
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
		//return emptyPiece
	} else if t.IsHazard() {
		return hazardPiece
	}

	return emptyPiece
}
