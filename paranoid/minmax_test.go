package paranoid

import (
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/moveset"
)

func TestOpponentMoves(t *testing.T) {
	// s _ _
	// h f e
	// _ _ s
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 2}},
	}
	s2 := api.Battlesnake{
		ID:     "s2",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, s2},
			Height: 3,
			Width:  3,
			Food: []api.Coord{
				{X: 1, Y: 1},
			},
		},
		You: me,
	}
	board := bitboard.CreateBitBoard(state)

	moves := withOppMoves(board, me.ID, moveset.SetDown(moveset.Create()))

	if len(moves) != 2 {
		panic("should be 2 sets of moves")
	}

}
