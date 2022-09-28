package mcts

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
)

func TestManhattanDistance(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ f _
	// _ _ h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 1, Y: 1}},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	if manhattanDistance(&board, 2, 4) != 2 {
		panic("distance wasn't 2!")
	}

	// f _ _
	// _ _ _
	// _ _ h
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state = api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 0, Y: 2}},
		},
		You: me,
	}
	board = b.BuildBoard(state)

	if manhattanDistance(&board, 2, 6) != 4 {
		panic("distance wasn't 4!")
	}
}

func TestFoodStrategy(t *testing.T) {
	//t.Skip()
	// _ _ f
	// _ _ _
	// _ _ h

	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 2, Y: 2}},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]

	move := foodStrategy(&board, id)

	if move.Dir != b.Up {
		panic("did not move toward food!")
	}

	// Moves toward closer food
	// f _ f
	// _ _ _
	// _ _ h

	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state = api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 2, Y: 2}, {X: 0, Y: 2}},
		},
		You: me,
	}
	board = b.BuildBoard(state)
	id = board.Ids["me"]

	move = foodStrategy(&board, id)

	if move.Dir != b.Up {
		panic("did not move toward food!")
	}
}

func TestCenterStrategy(t *testing.T) {
	//t.Skip()
	// _ _ f _ _
	// _ _ _ _ _
	// _ _ _ _ _
	// _ _ _ h _
	// _ _ _ _ _

	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 3, Y: 0},
		Body:   []api.Coord{{X: 3, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 0}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 5,
			Width:  5,
			Food:   []api.Coord{{X: 2, Y: 4}},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]

	move := centerStrategy(&board, id)

	if move.Dir != b.Left && move.Dir != b.Up {
		fmt.Println(move)
		panic("did not move toward Center!")
	}

}
