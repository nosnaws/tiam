package game

import (
	"fmt"
	"testing"
)

func TestVoronoi(t *testing.T) {
	// h _ f
	// s _ _
	// _ e e
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 2, Y: 2}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	//twoId := board.ids["two"]

	v := Voronoi(&board, id)

	if v.Score[id] != 1 {
		board.Print()
		fmt.Println(v)
		panic("Voronoi is should be 1!")
	}

	if v.FoodDepth[id] != 1 {
		board.Print()
		fmt.Println(v)
		panic("foodDepth should be 1!")
	}

	// f _ _ _ _
	// _ s s s h
	// s s s e f
	// _ _ _ _ _
	// _ _ _ _ _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 4, Y: 3}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}},
	}
	two = Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []Coord{{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food:   []Coord{{X: 4, Y: 2}, {X: 0, Y: 4}},
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.ids["me"]
	//twoId := board.ids["two"]

	v = Voronoi(&board, id)

	if v.Score[id] != 6 {
		board.Print()
		fmt.Println(v)
		panic("Voronoi is should be 5!")
	}

	if v.FoodDepth[id] != 0 {
		board.Print()
		fmt.Println(v)
		panic("Food depth is not 0")
	}

}
