package game

type Point struct {
	X int8
	Y int8
}

type Move string

const (
	Left  Move = "left"
	Right      = "right"
	Up         = "up"
	Down       = "down"
)

func indexToPoint(index uint16, width uint16) Point {
	return Point{X: int8(index % width), Y: int8(index / width)}
}

func pointToIndex(p Point, width uint16) uint16 {
	return uint16(int16(p.Y)*int16(width) + int16(p.X))
}

func indexInDirection(m Move, cur uint16, width, height uint16, isWrapped bool) uint16 {
	p := addPoints(indexToPoint(cur, width), moveToPoint(m))
	if isWrapped {
		p = adjustForWrapped(p, width, height)
	}
	return pointToIndex(p, width)
}

func adjustForWrapped(p Point, width, height uint16) Point {
	if p.X < 0 {
		return Point{X: int8(width - 1), Y: p.Y}
	}
	if p.Y < 0 {
		return Point{X: p.X, Y: int8(height - 1)}
	}
	if p.X >= int8(width) {
		return Point{X: 0, Y: p.Y}
	}
	if p.Y >= int8(height) {
		return Point{X: p.X, Y: 0}
	}
	return p
}

func moveToPoint(m Move) Point {
	if m == Left {
		return Point{X: -1, Y: 0}
	} else if m == Right {
		return Point{X: 1, Y: 0}
	} else if m == Up {
		return Point{X: 0, Y: 1}
	}
	return Point{X: 0, Y: -1}
}

func addPoints(a, b Point) Point {
	return Point{X: a.X + b.X, Y: a.Y + b.Y}
}

func GetCartesianProductOfMoves(board FastBoard) [][]SnakeMove {
	var allMoves [][]SnakeMove
	for id, h := range board.Healths {
		if h > 0 {
			moves := board.GetMovesForSnake(id)
			allMoves = append(allMoves, moves)
		}
	}

	var temp [][]SnakeMove
	for _, a := range allMoves[0] {
		temp = append(temp, []SnakeMove{a})
	}

	for i := 1; i < len(allMoves); i++ {
		temp = cartesianProduct(temp, allMoves[i])
	}

	return temp
}

func cartesianProduct(movesA [][]SnakeMove, movesB []SnakeMove) [][]SnakeMove {
	var result [][]SnakeMove
	for _, a := range movesA {
		for _, b := range movesB {
			var temp []SnakeMove
			for _, m := range a {
				temp = append(temp, m)
			}

			temp = append(temp, b)
			result = append(result, temp)
		}
	}

	return result
}
