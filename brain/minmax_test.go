package brain

import (
	"fmt"
	"math"
	"testing"

	"github.com/nosnaws/tiam/game"
)

func TestMinMaxDecision(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ s s s s h _ _ _ _ _
	// _ _ _ _ _ _ _ s s s h
	// _ _ _ _ _ _ s s s h f
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ s s h _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	me := game.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []game.Coord{{X: 10, Y: 5}, {X: 9, Y: 5}, {X: 8, Y: 5}, {X: 7, Y: 5}},
	}
	two := game.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []game.Coord{{X: 9, Y: 4}, {X: 8, Y: 4}, {X: 7, Y: 4}, {X: 6, Y: 4}},
	}
	three := game.Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []game.Coord{{X: 9, Y: 2}, {X: 8, Y: 2}, {X: 7, Y: 2}},
	}
	four := game.Battlesnake{
		ID:     "four",
		Health: 100,
		Body:   []game.Coord{{X: 5, Y: 6}, {X: 4, Y: 6}, {X: 3, Y: 6}, {X: 2, Y: 6}, {X: 1, Y: 6}},
	}
	state := game.GameState{
		Turn: 0,
		Board: game.Board{
			Snakes: []game.Battlesnake{me, two, three, four},
			Height: 11,
			Width:  11,
			Food:   []game.Coord{{X: 10, Y: 4}},
		},
		You: me,
	}
	board := game.BuildBoard(state)

	move := AlphaBeta(board, 6, int32(math.Inf(-1)), int32(math.Inf(1)), true, game.MeId, game.SnakeMove{}).Move
	if move.Dir != game.Up {
		board.Print()
		fmt.Println("selected ", move.Dir)
		panic("Did not go up!")
	}
}
