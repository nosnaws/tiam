package mcts

import (
	"math"
	"math/rand"

	"github.com/nosnaws/tiam/board"
)

func mixedStrategy(fb *board.FastBoard, id board.SnakeId) board.SnakeMove {
	r := rand.Float32()

	if r > 0.9 {
		return defensiveStrategy(fb, id)
	} else if r > 0.8 {
		return aggressiveStrategy(fb, id)
	} else if r > 0.6 {
		return foodStrategy(fb, id)
	} else {
		return randomStrategy(fb, id)
	}
}

func randomStrategy(fb *board.FastBoard, id board.SnakeId) board.SnakeMove {
	sMoves := fb.GetMovesForSnake(id)
	randomMove := sMoves[rand.Intn(len(sMoves))]

	return randomMove
}

func foodStrategy(fb *board.FastBoard, id board.SnakeId) board.SnakeMove {
	moves := fb.GetMovesForSnake(id)
	food := findAll(fb, func(t board.Tile) bool {
		if t.IsFood() {
			return true
		}
		return false
	})

	if len(food) < 1 {
		return randomStrategy(fb, id)
	}

	// best move is the closest to food
	nearest := float64(-1)
	var nMove board.SnakeMove

	for _, foodI := range food {
		for _, move := range moves {
			mIndex := fb.MoveToIndex(move)
			if nearest == -1 {
				nearest = manhattanDistance(fb, mIndex, foodI)
				nMove = move
				continue
			}

			dis := manhattanDistance(fb, mIndex, foodI)
			if dis < nearest {
				nearest = dis
				nMove = move
			}
		}
	}

	return nMove
}

func centerStrategy(fb *board.FastBoard, id board.SnakeId) board.SnakeMove {
	moves := fb.GetMovesForSnake(id)
	centerIndex := uint16(len(fb.List) / 2)

	// best move is the closest to the center
	nearest := float64(-1)
	var nMove board.SnakeMove

	for _, move := range moves {
		mIndex := fb.MoveToIndex(move)
		if nearest == -1 {
			nearest = manhattanDistance(fb, mIndex, centerIndex)
			nMove = move
			continue
		}

		dis := manhattanDistance(fb, mIndex, centerIndex)
		if dis < nearest {
			nearest = dis
			nMove = move
		}
	}

	return nMove
}

func defensiveStrategy(fb *board.FastBoard, id board.SnakeId) board.SnakeMove {
	moves := fb.GetMovesForSnake(id)

	nearestSnakeId := findNearestSnake(fb, id)

	if nearestSnakeId == 0 {
		return randomStrategy(fb, id)
	}

	// best move is the farthest from the snake
	farthest := float64(-1)
	var nMove board.SnakeMove

	for _, move := range moves {
		mIndex := fb.MoveToIndex(move)
		dis := manhattanDistance(fb, mIndex, fb.Heads[nearestSnakeId])
		if dis > farthest || farthest == -1 {
			farthest = dis
			nMove = move
		}
	}

	return nMove
}

func aggressiveStrategy(fb *board.FastBoard, id board.SnakeId) board.SnakeMove {
	moves := fb.GetMovesForSnake(id)

	nearestSnakeId := findNearestSnake(fb, id)

	if nearestSnakeId == 0 {
		return randomStrategy(fb, id)
	}

	// best move is the closest to the snake
	nearest := float64(-1)
	var nMove board.SnakeMove

	for _, move := range moves {
		mIndex := fb.MoveToIndex(move)
		dis := manhattanDistance(fb, mIndex, fb.Heads[nearestSnakeId])
		if dis < nearest || nearest == -1 {
			nearest = dis
			nMove = move
		}
	}

	return nMove
}

func findAll(fb *board.FastBoard, cb func(board.Tile) bool) []uint16 {
	found := []uint16{}
	for i, tile := range fb.List {
		if cb(tile) {
			found = append(found, uint16(i))
		}
	}

	return found
}

func countAliveSnakes(fb *board.FastBoard) int {
	num := 0

	for id := range fb.Heads {
		if fb.IsSnakeAlive(id) {
			num += 1
		}
	}
	return num
}

func findNearestSnake(fb *board.FastBoard, id board.SnakeId) board.SnakeId {
	myHead := fb.Heads[id]
	var nearest uint16
	var nearestId board.SnakeId
	for sId, head := range fb.Heads {
		closestDis := manhattanDistance(fb, myHead, nearest)
		currentDis := manhattanDistance(fb, myHead, head)
		if sId != id && (currentDis < closestDis || nearestId == 0) {
			nearest = head
			nearestId = sId
		}
	}

	return nearestId
}

func getLongestSnakeId(fb *board.FastBoard) board.SnakeId {
	lId := board.SnakeId(0)
	lLength := uint8(0)
	for id, length := range fb.Lengths {
		if fb.IsSnakeAlive(id) {
			if lLength == 0 {
				lId = id
				lLength = length
				continue
			}

			if length > lLength {
				lId = id
				lLength = length
			}
		}
	}

	return lId
}

func manhattanDistance(fb *board.FastBoard, a, b uint16) float64 {
	if fb.IsWrapped {
		return manhattanDistanceWrapped(a, b, fb.Width, fb.Height)
	}

	return manhattanDistanceStandard(a, b, fb.Width)
}

func manhattanDistanceStandard(a, b, width uint16) float64 {
	aP, bP := board.IndexToPoint(a, width), board.IndexToPoint(b, width)

	return calcDistance(aP, bP)
}

func manhattanDistanceWrapped(a, b, width, height uint16) float64 {
	aP, bP := board.IndexToPoint(a, width), board.IndexToPoint(b, width)
	leftTranspose := board.Point{X: bP.X - int8(width), Y: bP.Y}
	rightTranspose := board.Point{X: bP.X + int8(width), Y: bP.Y}
	upTranspose := board.Point{Y: bP.Y + int8(height), X: bP.X}
	downTranspose := board.Point{Y: bP.Y - int8(height), X: bP.X}

	// start with non-wrapped distance
	currentBest := calcDistance(aP, bP)

	leftWDis := calcDistance(aP, leftTranspose)
	if leftWDis < currentBest {
		currentBest = leftWDis
	}

	rightWDis := calcDistance(aP, rightTranspose)
	if rightWDis < currentBest {
		currentBest = rightWDis
	}

	upWDis := calcDistance(aP, upTranspose)
	if upWDis < currentBest {
		currentBest = upWDis
	}

	downWDis := calcDistance(aP, downTranspose)
	if downWDis < currentBest {
		currentBest = downWDis
	}

	return currentBest
}

func calcDistance(pointA, pointB board.Point) float64 {
	return math.Abs(float64(pointA.X)-float64(pointB.X)) + math.Abs(float64(pointA.Y)-float64(pointB.Y))
}
