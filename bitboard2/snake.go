package bitboard

import (
	"fmt"
	num "github.com/shabbyrobe/go-num"

	api "github.com/nosnaws/tiam/battlesnake"
)

// [0] = head, [1] = tail, [2] = segment after tail
type headtail [3]int

type snake struct {
	health        int
	body          []int
	board         num.U128
	headBoard     num.U128
	hasDoubleTail bool
	Length        int
}

func createSnake(s api.Battlesnake, width int) *snake {
	board := num.U128From16(0)
	headBoard := num.U128From16(0)

	body := []int{}

	for b := len(s.Body) - 1; b >= 0; b-- {
		i := getIndex(s.Body[b], width)
		board = board.SetBit(i, 1)
		body = append(body, i)
	}

	head := getIndex(s.Body[0], width)
	tail := getIndex(s.Body[len(s.Body)-1], width)
	afterTail := getIndex(s.Body[len(s.Body)-2], width)

	headBoard = headBoard.SetBit(head, 1)

	return &snake{
		health:        int(s.Health),
		body:          body,
		Length:        len(body),
		board:         board,
		headBoard:     headBoard,
		hasDoubleTail: afterTail == tail,
	}
}

// TODO: keep this in state so we don't have to recreate it each call
func (s *snake) getHeadBoard() num.U128 {
	return s.headBoard
	//b := big.NewInt(0)
	//b.SetBit(b, s.GetHeadIndex(), 1)

	//return b
}

// second turn or ate last turn
func (s *snake) stackedTail() bool {
	// should never happen
	if s.Length < 2 {
		return false
	}

	return s.body[0] == s.body[1]
}

func (s *snake) GetHeadIndex() int {
	return s.body[len(s.body)-1]
}

func (s *snake) getNeckIndex() int {
	return s.body[len(s.body)-2]
}

func (s *snake) setHeadIndex(i int) {
	s.body = append(s.body, i)
}

func (s *snake) getTailIndex() int {
	return s.body[0]
}

func (s *snake) GetHealth() int {
	return s.health
}

func (s *snake) moveHead(newIdx int, dir Dir, width uint) {
	if dir == Left {
		s.headBoard = s.headBoard.Rsh(1)
	} else if dir == Right {
		s.headBoard = s.headBoard.Lsh(1)
	} else if dir == Up {
		s.headBoard = s.headBoard.Lsh(width)
	} else {
		s.headBoard = s.headBoard.Rsh(width)
	}
	//s.headBoard = s.headBoard.SetBit(s.GetHeadIndex(), 0)
	//s.headBoard = s.headBoard.SetBit(i, 1)

	s.board = s.board.Or(s.headBoard)

	s.setHeadIndex(newIdx)
}

func (s *snake) moveTail() {
	isTripleStacked := s.body[0] == s.body[1] && s.body[0] == s.body[2]

	// move tail out of location if it isn't stacked
	if !isTripleStacked && !s.stackedTail() {
		s.board = s.board.SetBit(s.getTailIndex(), 0)
	}

	s.body = s.body[1:]
}

func (s *snake) feed() {
	s.health = 100
	tail := s.getTailIndex()
	s.body = append([]int{tail}, s.body...)
	s.Length += 1
}

func (s *snake) IsAlive() bool {
	return s.health > 0
}

func (s *snake) kill() {
	s.health = 0
	s.body = []int{}
	s.board = num.U128From16(0)
	s.headBoard = num.U128From16(0)
	s.Length = 0
}

func (s *snake) clone() *snake {
	body := make([]int, len(s.body))
	copy(body, s.body)

	board := s.board
	headBoard := s.headBoard

	return &snake{
		health:        s.health,
		Length:        s.Length,
		hasDoubleTail: s.hasDoubleTail,
		body:          body,
		board:         board,
		headBoard:     headBoard,
	}
}

func (s *snake) print() {
	fmt.Println("health:", s.health)
	fmt.Println("length:", s.Length)
	fmt.Println("body:", s.body)
}
