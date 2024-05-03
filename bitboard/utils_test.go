package bitboard

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
)

func TestCartesianProduct(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ e _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ e _ _ _ _ e _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _

	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 8},
		Body: []api.Coord{
			{X: 2, Y: 8},
			{X: 2, Y: 8},
			{X: 2, Y: 8},
		},
	}
	s1 := api.Battlesnake{
		ID:     "s1",
		Health: 100,
		Head:   api.Coord{X: 7, Y: 8},
		Body: []api.Coord{
			{X: 7, Y: 8},
			{X: 7, Y: 8},
			{X: 7, Y: 8},
		},
	}
	s2 := api.Battlesnake{
		ID:     "s2",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 3},
		Body: []api.Coord{
			{X: 2, Y: 3},
			{X: 2, Y: 3},
			{X: 2, Y: 3},
		},
	}
	s3 := api.Battlesnake{
		ID:     "s3",
		Health: 100,
		Head:   api.Coord{X: 7, Y: 3},
		Body: []api.Coord{
			{X: 7, Y: 3},
			{X: 7, Y: 3},
			{X: 7, Y: 3},
		},
	}

	state := api.GameState{
		Turn: 1,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, s1, s2, s3},
			Height: 11,
			Width:  11,
		},
		You: me,
	}

	board := CreateBitBoard(state)

	states := board.GetCartesianProductOfMoves()
	if len(states) != 256 {
		board.Print()
		fmt.Println(len(states))
		panic("Did not get correct states!")
	}
}
