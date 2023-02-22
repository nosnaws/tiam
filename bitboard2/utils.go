package bitboard

import (
	"fmt"

	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/nosnaws/tiam/moveset"
	"github.com/shabbyrobe/go-num"
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

func (bb *BitBoard) GetOpponents() []*snake {
	snakes := []*snake{}

	for id, snake := range bb.Snakes {
		if id != bb.meId && snake.IsAlive() {
			snakes = append(snakes, snake)
		}
	}

	return snakes
}

func (bb *BitBoard) indexToString(i int) string {
	board := num.U128From16(0)
	board = board.SetBit(i, 1)

	for id, s := range bb.Snakes {
		if !s.IsAlive() {
			continue
		}

		if s.GetHeadIndex() == i {
			return fmt.Sprintf(" %sh ", firstN(id, 1))
		}
		if board.And(s.board).BitLen() > 0 {
			return " ss "
		}
	}

	if board.And(bb.food).BitLen() > 0 {
		return " ff "
	}

	if board.And(bb.hazards).BitLen() > 0 {
		return " zz "
	}

	//test = big.NewInt(0)
	//if test.And(board, bb.empty).BitLen() > 0 {
	//return " ee "
	//}

	return " __ "
}

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func (bb *BitBoard) Print() {
	fmt.Println("######")
	for id, s := range bb.Snakes {
		fmt.Printf("%s - length:%d\n", id, s.Length)
	}
	for id, s := range bb.Snakes {
		fmt.Printf("%s - health:%d\n", id, s.health)
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

func (bb *BitBoard) printBoard(b num.U128) {
	for y := int(bb.height - 1); y >= 0; y-- {
		var line string
		for x := 0; x < int(bb.width); x++ {
			p := api.Coord{X: x, Y: y}
			index := getIndex(p, bb.width)
			line = line + fmt.Sprint(b.Bit(index))
		}
		fmt.Println(line)
	}
}

func SplitSnakeMoveSet(moves SnakeMoveSet) []SnakeMoveSet {
	split := []SnakeMoveSet{}
	for _, move := range moveset.Split(moves.Set) {
		split = append(split, SnakeMoveSet{Id: moves.Id, Set: move})
	}

	return split
}

func (bb *BitBoard) GetCartesianProductOfMoves() [][]SnakeMoveSet {
	var allMoves [][]SnakeMoveSet
	for id, snake := range bb.Snakes {
		if snake.IsAlive() {
			moves := SplitSnakeMoveSet(bb.GetMoves(id))
			allMoves = append(allMoves, moves)
		}
	}

	var temp [][]SnakeMoveSet
	for _, a := range allMoves[0] {
		temp = append(temp, []SnakeMoveSet{a})
	}

	for i := 1; i < len(allMoves); i++ {
		temp = CartesianProduct(temp, allMoves[i])
	}

	return temp
}

func CartesianProduct(movesA [][]SnakeMoveSet, movesB []SnakeMoveSet) [][]SnakeMoveSet {
	var result [][]SnakeMoveSet
	for _, a := range movesA {
		for _, b := range movesB {
			var temp []SnakeMoveSet
			for _, m := range a {
				temp = append(temp, m)
			}

			temp = append(temp, b)
			result = append(result, temp)
		}
	}

	return result
}

// updates without removing snakes
func (bb *BitBoard) Update(new *BitBoard) {
	//bb.empty = new.empty
	bb.food = new.food
	bb.hazards = new.hazards
}

func (bb *BitBoard) GetLastSnakeMoveFromExternal(snake api.Battlesnake) Dir {
	head := getIndex(snake.Body[0], bb.width)
	neck := getIndex(snake.Body[1], bb.width)

	if indexInDirection(Left, neck, bb.width, bb.height, bb.isWrapped) == head {
		return Left
	}
	if indexInDirection(Right, neck, bb.width, bb.height, bb.isWrapped) == head {
		return Right
	}
	if indexInDirection(Up, neck, bb.width, bb.height, bb.isWrapped) == head {
		return Up
	}

	return Down
}
