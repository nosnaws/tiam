package mmm

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
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

func TestVoronoiScore(t *testing.T) {
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

	board := b.BuildBoard(state)
	id := board.Ids["me"]
	twoId := board.Ids["two"]

	score := StrategyV4(&board, id, twoId, 1)

	if true {
		board.Print()
		fmt.Println("SCORE", score)
		panic("bad voronoi")
	}
}
