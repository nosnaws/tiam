package board

import (
	"fmt"
	api "github.com/nosnaws/tiam/battlesnake"
	"testing"
)

func CreateIslandBridgesGame() api.GameState {

	islandBridgesHz := []api.Coord{
		{X: 0, Y: 1},
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 4, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 1},
		{X: 6, Y: 0},
		{X: 9, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 1},
		{X: 10, Y: 4},
		{X: 9, Y: 5},
		{X: 10, Y: 5},
		{X: 10, Y: 6},
		{X: 10, Y: 9},
		{X: 10, Y: 10},
		{X: 9, Y: 10},
		{X: 6, Y: 10},
		{X: 5, Y: 10},
		{X: 5, Y: 9},
		{X: 4, Y: 10},
		{X: 1, Y: 10},
		{X: 0, Y: 10},
		{X: 0, Y: 9},
		{X: 0, Y: 6},
		{X: 0, Y: 5},
		{X: 1, Y: 5},
		{X: 0, Y: 4},
		{X: 2, Y: 4},
		{X: 3, Y: 5},
		{X: 4, Y: 5},
		{X: 5, Y: 5},
		{X: 6, Y: 5},
		{X: 7, Y: 5},
		{X: 5, Y: 7},
		{X: 5, Y: 6},
		{X: 5, Y: 5},
		{X: 5, Y: 4},
		{X: 5, Y: 3},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Height:  11,
			Width:   11,
			Hazards: islandBridgesHz,
		},
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
				Settings: api.Settings{
					HazardDamagePerTurn: 100,
				},
			},
		},
	}

	return state
}

func TestVoronoi(t *testing.T) {
	t.Skip()
	// h _ f
	// s _ _
	// _ e e
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
	}
	state := api.GameState{
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 2, Y: 2}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]
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
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 4, Y: 3}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}},
	}
	two = api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state = api.GameState{
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food:   []api.Coord{{X: 4, Y: 2}, {X: 0, Y: 4}},
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.Ids["me"]
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

func TestVoronoiScore(t *testing.T) {
	t.Skip()
	// z z h e z z z _ _ z z
	// z e e e _ z _ _ _ _ z
	// _ e e s s s s _ _ _ _
	// _ _ e s f z _ _ _ _ _
	// z _ h s _ z _ _ _ _ z
	// z z _ z z z z z f z z
	// z _ _ _ _ z _ _ _ _ z
	// _ _ _ _ _ z _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// z _ _ _ _ z _ _ _ _ z
	// z z _ _ z z z _ _ z z

	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []api.Coord{
			{X: 2, Y: 6},
			{X: 3, Y: 6},
			{X: 3, Y: 7},
			{X: 3, Y: 8},
			{X: 4, Y: 8},
			{X: 5, Y: 8},
			{X: 6, Y: 8},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{X: 2, Y: 10},
			{X: 3, Y: 10},
			{X: 3, Y: 9},
			{X: 2, Y: 9},
			{X: 1, Y: 9},
			{X: 1, Y: 8},
			{X: 2, Y: 8},
			{X: 2, Y: 7},
		},
	}

	state := CreateIslandBridgesGame()
	state.Board.Food = []api.Coord{
		{X: 4, Y: 7},
		{X: 8, Y: 5},
	}
	state.Board.Snakes = []api.Battlesnake{me, two}
	state.You = me

	board := BuildBoard(state)
	id := board.Ids["me"]
	//twoId := board.ids["two"]

	v := Voronoi(&board, id)

	if true {
		board.Print()
		fmt.Println("FOOD", v.FoodDepth)
		fmt.Println("V", v.Score)
		panic("bad voronoi")
	}
}
