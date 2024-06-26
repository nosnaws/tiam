package board

import b "github.com/nosnaws/tiam/battlesnake"

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

var leftPoint = Point{X: -1, Y: 0}
var rightPoint = Point{X: 1, Y: 0}
var upPoint = Point{X: 0, Y: 1}
var downPoint = Point{X: 0, Y: -1}

func IndexToPoint(index uint16, width uint16) Point {
	return Point{X: int8(index % width), Y: int8(index / width)}
}

func pointToIndex(p Point, width uint16) uint16 {
	return uint16(int16(p.Y)*int16(width) + int16(p.X))
}

func pointInDirection(m Move, cur uint16, width, height uint16, isWrapped bool) Point {
	p := addPoints(IndexToPoint(cur, width), moveToPoint(m))
	if isWrapped {
		p = adjustForWrapped(p, width, height)
	}
	return p
}

func IndexInDirection(m Move, cur uint16, width, height uint16, isWrapped bool) uint16 {
	p := pointInDirection(m, cur, width, height, isWrapped)
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
		return leftPoint
	} else if m == Right {
		return rightPoint
	} else if m == Up {
		return upPoint
	}
	return downPoint
}

func addPoints(a, b Point) Point {
	return Point{X: a.X + b.X, Y: a.Y + b.Y}
}

func isSnakeInitialState(s b.Battlesnake) bool {
	head := s.Body[0]
	bp1 := s.Body[1]
	bp2 := s.Body[2]

	return head == bp1 && head == bp2
}

func isSnakeDoubleButt(s b.Battlesnake) bool {
	l := len(s.Body)
	tail := s.Body[l-1]
	beforeTail := s.Body[l-2]

	return tail == beforeTail
}

func GetCartesianProductOfMoves(board *FastBoard) [][]SnakeMove {
	var allMoves [][]SnakeMove
	for id := range board.Healths {
		if board.IsSnakeAlive(id) {
			moves := board.GetMovesForSnake(id)
			allMoves = append(allMoves, moves)
		}
	}

	var temp [][]SnakeMove
	for _, a := range allMoves[0] {
		temp = append(temp, []SnakeMove{a})
	}

	for i := 1; i < len(allMoves); i++ {
		temp = CartesianProduct(temp, allMoves[i])
	}

	return temp
}

func CartesianProduct(movesA [][]SnakeMove, movesB []SnakeMove) [][]SnakeMove {
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

func MovesToMap(moves []SnakeMove) map[SnakeId]Move {
	m := make(map[SnakeId]Move, len(moves))
	for _, move := range moves {
		m[move.Id] = move.Dir
	}
	return m
}
