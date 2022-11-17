package bitboard

import (
	"fmt"
	"math/big"

	api "github.com/nosnaws/tiam/battlesnake"
)

type Dir string

const (
	Left  Dir = "left"
	Right     = "right"
	Up        = "up"
	Down      = "down"
)

var leftPoint = api.Coord{X: -1, Y: 0}
var rightPoint = api.Coord{X: 1, Y: 0}
var upPoint = api.Coord{X: 0, Y: 1}
var downPoint = api.Coord{X: 0, Y: -1}

func coordInDirection(m Dir, cur, width, height int, isWrapped bool) api.Coord {
	p := addCoords(indexToPoint(cur, width), moveToCoord(m))
	if isWrapped {
		p = adjustForWrapped(p, width, height)
	}
	return p
}

func indexInDirection(m Dir, cur, width, height int, isWrapped bool) int {
	p := coordInDirection(m, cur, width, height, isWrapped)
	return getIndex(p, width)
}

func isDirOutOfBounds(m Dir, cur, width, height int, isWrapped bool) bool {
	c := coordInDirection(m, cur, width, height, isWrapped)

	if c.X < 0 || c.X >= width {
		return true
	}
	if c.Y < 0 || c.Y >= height {
		return true
	}
	return false
}

func adjustForWrapped(p api.Coord, width, height int) api.Coord {
	if p.X < 0 {
		return api.Coord{X: width - 1, Y: p.Y}
	}
	if p.Y < 0 {
		return api.Coord{X: p.X, Y: height - 1}
	}
	if p.X >= width {
		return api.Coord{X: 0, Y: p.Y}
	}
	if p.Y >= height {
		return api.Coord{X: p.X, Y: 0}
	}
	return p
}

func moveToCoord(m Dir) api.Coord {
	if m == Left {
		return leftPoint
	} else if m == Right {
		return rightPoint
	} else if m == Up {
		return upPoint
	}
	return downPoint
}

func indexToPoint(index, width int) api.Coord {
	return api.Coord{X: int(index % width), Y: int(index / width)}
}

func addCoords(a, b api.Coord) api.Coord {
	return api.Coord{X: a.X + b.X, Y: a.Y + b.Y}
}

func getIndex(p api.Coord, width int) int {
	return p.Y*width + p.X
}

func (bb *BitBoard) indexToString(i int) string {
	board := big.NewInt(0)
	board.SetBit(board, i, 1)

	for id, s := range bb.snakes {
		if s.getHeadIndex() == i {
			return fmt.Sprintf(" %dh ", id)
		}
		test := big.NewInt(0)
		if test.And(board, s.board).BitLen() > 0 {
			return " ss "
		}
	}

	test := big.NewInt(0)
	if test.And(board, bb.food).BitLen() > 0 {
		return " ff "
	}

	test = big.NewInt(0)
	if test.And(board, bb.hazards).BitLen() > 0 {
		return " zz "
	}

	//test = big.NewInt(0)
	//if test.And(board, bb.empty).BitLen() > 0 {
	//return " ee "
	//}

	return " __ "
}

func (bb *BitBoard) Print() {
	fmt.Println("######")
	for id, s := range bb.snakes {
		fmt.Printf("%d - length:%d\n", id, s.length)
	}
	for id, s := range bb.snakes {
		fmt.Printf("%d - health:%d\n", id, s.health)
	}

	for y := int(bb.height - 1); y >= 0; y-- {
		var line string
		for x := 0; x < int(bb.width); x++ {
			p := api.Coord{X: x, Y: y}
			line = line + bb.indexToString(getIndex(p, bb.width))
		}
		fmt.Println(line)
	}
}
