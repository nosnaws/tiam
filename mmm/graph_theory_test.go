package mmm

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
)

func TestArticulationPoints(t *testing.T) {
	//t.Skip()
	// z z _ _ z z z _ _ z z
	// z s s _ _ z _ _ _ _ z
	// _ _ s _ _ _ _ _ _ h _
	// _ _ e _ _ z s s s s _
	// z _ _ _ _ z s s _ _ z
	// z z _ _ z z z _ _ z z
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []api.Coord{
			{9, 3},
			{9, 2},
			{8, 2},
			{7, 2},
			{7, 1},
			{6, 1},
			{6, 2},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{2, 2},
			{2, 3},
			{2, 4},
			{1, 4},
		},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 6,
			Width:  11,
			Hazards: []api.Coord{
				{X: 0, Y: 4},
				{X: 0, Y: 5},
				{X: 1, Y: 5},
				{X: 4, Y: 5},
				{X: 5, Y: 5},
				{X: 5, Y: 4},
				{X: 6, Y: 5},
				{X: 9, Y: 5},
				{X: 10, Y: 5},
				{X: 10, Y: 4},
				{X: 10, Y: 1},
				{X: 10, Y: 0},
				{X: 9, Y: 0},
				{X: 6, Y: 0},
				{X: 5, Y: 0},
				{X: 5, Y: 1},
				{X: 5, Y: 2},
				{X: 4, Y: 0},
				{X: 1, Y: 0},
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
		},
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
				Settings: api.Settings{
					HazardDamagePerTurn: 100,
				},
			},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	//id := board.Ids["me"]

	arts := getArticulationPoints(&board, 42)

	if true {
		board.Print()
		fmt.Println(arts)
		panic("wrong art points !")
	}

}
