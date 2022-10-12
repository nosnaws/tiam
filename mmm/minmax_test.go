package mmm

import (
	"context"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/nosnaws/tiam/board"
)

func TestDoesNotLoopForever(t *testing.T) {
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ s s s _ _ e s _
	// _ _ z z z s z z z s _
	// _ _ z _ _ s _ _ z s _
	// _ _ z _ _ e _ _ z s _
	// _ _ f _ _ _ _ _ _ s _
	// _ s z _ _ _ _ _ z _ _
	// _ s z _ _ _ _ _ z _ _
	// _ s z z z f z z z _ _
	// _ s _ _ _ _ _ _ _ _ e
	// _ h _ _ _ _ _ _ _ s s
	me := api.Battlesnake{
		ID:     "me",
		Health: 20,
		Body: []api.Coord{
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 1, Y: 2},
			{X: 1, Y: 3},
			{X: 1, Y: 4},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{X: 5, Y: 6},
			{X: 5, Y: 7},
			{X: 5, Y: 8},
			{X: 5, Y: 9},
			{X: 4, Y: 9},
			{X: 3, Y: 9},
		},
	}
	three := api.Battlesnake{
		ID:     "three",
		Health: 100,
		Body: []api.Coord{
			{X: 8, Y: 9},
			{X: 9, Y: 9},
			{X: 9, Y: 8},
			{X: 9, Y: 7},
			{X: 9, Y: 6},
			{X: 9, Y: 5},
		},
	}
	four := api.Battlesnake{
		ID:     "four",
		Health: 40,
		Body: []api.Coord{
			{X: 10, Y: 1},
			{X: 10, Y: 0},
			{X: 9, Y: 0},
		},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two, three, four},
			Height: 11,
			Width:  11,
			Food: []api.Coord{
				{X: 2, Y: 5},
				{X: 5, Y: 2},
			},
			Hazards: []api.Coord{
				{X: 2, Y: 6},
				{X: 2, Y: 6},
				{X: 2, Y: 7},
				{X: 2, Y: 8},
				{X: 3, Y: 8},
				{X: 4, Y: 8},
				{X: 4, Y: 8},
				{X: 6, Y: 8},
				{X: 6, Y: 8},
				{X: 7, Y: 8},
				{X: 8, Y: 8},
				{X: 8, Y: 7},
				{X: 8, Y: 6},
				{X: 8, Y: 6},
				{X: 8, Y: 4},
				{X: 8, Y: 4},
				{X: 8, Y: 3},
				{X: 8, Y: 2},
				{X: 7, Y: 2},
				{X: 6, Y: 2},
				{X: 6, Y: 2},
				{X: 4, Y: 2},
				{X: 4, Y: 2},
				{X: 3, Y: 2},
				{X: 2, Y: 2},
				{X: 2, Y: 3},
				{X: 2, Y: 4},
				{X: 2, Y: 4},
			},
		},
		Game: api.Game{
			Ruleset: api.Ruleset{
				Settings: api.Settings{
					HazardDamagePerTurn: 50,
				},
			},
		},
		You: me,
	}
	gs := board.BuildBoard(state)
	cache := CreateCache(&gs, 0)

	move, _, _ := MultiMinmax(context.TODO(), cache, &gs, StrategyV3, 12, "")

	if move != board.Left && move != board.Right {
		panic("Did not go left or right!")
	}
}
