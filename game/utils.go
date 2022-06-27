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
	return uint16(p.Y)*width + uint16(p.X)
}

func indexInDirection(m Move, cur uint16, width uint16) uint16 {
	return pointToIndex(addPoints(indexToPoint(cur, width), moveToPoint(m)), width)
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
