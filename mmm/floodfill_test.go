package mmm

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
)

func TestFloodfill(t *testing.T) {
	//t.Skip()
	// s s s _ _
	// s _ s _ _
	// s h s s x
	// _ _ _ e x
	// _ _ _ x x
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []api.Coord{
			{1, 2},
			{0, 2},
			{0, 3},
			{0, 4},
			{1, 4},
			{2, 4},
			{2, 3},
			{2, 2},
			{2, 2},
			{3, 2},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{3, 1},
			{3, 0},
			{4, 0},
			{4, 1},
			{4, 2},
		},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 11,
			Width:  11,
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]

	ff, _, _ := floodfill(&board, int(board.Heads[id]), 11, []uint16{})

	if ff != 8 {
		board.Print()
		fmt.Println("floodfill score", ff)
		panic("wrong floodfill amount!")
	}

}

func TestFloodfillFood(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// s _ _ _ _
	// s h _ _ _
	// _ f _ e _
	// _ _ _ s s
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []api.Coord{
			{1, 2},
			{0, 2},
			{0, 3},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{3, 1},
			{3, 0},
			{4, 0},
		},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food: []api.Coord{
				{1, 1},
			},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]

	ff, foodDepth, _ := floodfill(&board, int(board.Heads[id]), 11, []uint16{})

	if foodDepth != 1 {
		board.Print()
		fmt.Println("food depth", ff)
		panic("wrong food depth!")
	}

}
