package bitboard

import (
	"fmt"
	"math/big"

	api "github.com/nosnaws/tiam/battlesnake"
)

// [0] = head, [1] = tail, [2] = segment after tail
type headtail [3]int

type snake struct {
	health        int
	body          []int
	board         *big.Int
	headBoard     *big.Int
	hasDoubleTail bool
	Length        int
}

func createSnake(s api.Battlesnake, width int) *snake {
	board := big.NewInt(0)
	headBoard := big.NewInt(0)

	body := []int{}

	for b := len(s.Body) - 1; b >= 0; b-- {
		i := getIndex(s.Body[b], width)
		board.SetBit(board, i, 1)
		body = append(body, i)
	}

	head := getIndex(s.Body[0], width)
	tail := getIndex(s.Body[len(s.Body)-1], width)
	afterTail := getIndex(s.Body[len(s.Body)-2], width)

	headBoard.SetBit(headBoard, head, 1)

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
func (s *snake) getHeadBoard() *big.Int {
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

func (s *snake) moveHead(i int) {
	s.headBoard.SetBit(s.headBoard, s.GetHeadIndex(), 0)
	s.headBoard.SetBit(s.headBoard, i, 1)

	s.board.SetBit(s.board, i, 1)
	s.setHeadIndex(i)
}

func (s *snake) moveTail() {
	isTripleStacked := s.body[0] == s.body[1] && s.body[0] == s.body[2]

	// move tail out of location if it isn't stacked
	if !isTripleStacked && !s.stackedTail() {
		s.board.SetBit(s.board, s.getTailIndex(), 0)
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
	s.board = nil
	s.headBoard = nil
	s.Length = 0
}

func (s *snake) clone() *snake {
	body := make([]int, len(s.body))
	copy(body, s.body)

	board := big.NewInt(0).Set(s.board)
	headBoard := big.NewInt(0).Set(s.headBoard)

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
