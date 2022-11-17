package bitboard

import (
	"math/big"
	"sort"

	api "github.com/nosnaws/tiam/battlesnake"
)

// [0] = head, [1] = tail, [2] = segment after tail
type headtail [3]int

type snake struct {
	health        int
	body          []int
	board         *big.Int
	hasDoubleTail bool
	length        int
}

func createSnake(s api.Battlesnake, width int) *snake {
	board := big.NewInt(0)

	body := []int{}
	for _, bodyPart := range s.Body {
		i := getIndex(bodyPart, width)
		board.SetBit(board, i, 1)
		body = append(body, i)
	}

	// allows for push/pop in move operations
	sort.Sort(sort.Reverse(sort.IntSlice(body)))

	tail := getIndex(s.Body[len(s.Body)-1], width)
	afterTail := getIndex(s.Body[len(s.Body)-2], width)

	return &snake{
		health:        int(s.Health),
		body:          body,
		length:        int(s.Length),
		board:         board,
		hasDoubleTail: afterTail == tail,
	}
}

// second turn or ate last turn
func (s *snake) stackedTail() bool {
	// should never happen
	if s.length < 2 {
		return false
	}

	return s.body[0] == s.body[1]
}

func (s *snake) getHeadIndex() int {
	return s.body[s.length-1]
}

func (s *snake) setHeadIndex(i int) {
	s.body = append(s.body, i)
}

func (s *snake) getTailIndex() int {
	return s.body[0]
}

func (s *snake) moveHead(i int) {
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
	s.body = append(s.body, s.body[0])
}

func (s *snake) isAlive() bool {
	return s.length > 0
}

func (s *snake) kill() {
	s.health = 0
	s.body = []int{}
	s.board.SetInt64(0)
	s.length = 0
}
