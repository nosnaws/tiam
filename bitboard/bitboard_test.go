package bitboard

import (
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
)

func TestCreateBitBoard(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// _ s s s h
	// s s s e f
	// _ _ _ _ _
	// _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 4, Y: 3}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state := api.GameState{
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food:   []api.Coord{{X: 4, Y: 2}},
		},
		You: me,
	}

	bb := CreateBitBoard(state)
	bb.Print()
}
