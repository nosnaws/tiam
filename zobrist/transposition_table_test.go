package zobrist

import (
	"testing"

	g "github.com/nosnaws/tiam/game"
)

func TestTTable(t *testing.T) {

	// _ _ _ _ _
	// _ s s s h
	// s s s e f
	// _ _ _ _ _
	// _ _ _ _ _
	me := g.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []g.Coord{{X: 4, Y: 3}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}},
	}
	two := g.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []g.Coord{{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state := g.GameState{
		Board: g.Board{
			Snakes: []g.Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food:   []g.Coord{{X: 4, Y: 2}},
		},
		You: me,
	}
	board := g.BuildBoard(state)

	tt := InitializeZobristHash(5, 5)

	hash1 := GetZobristKey(tt, board)
	hash2 := GetZobristKey(tt, board)

	if hash1 != hash2 {
		panic("Hashes did not match!")
	}

	moves := make(map[g.SnakeId]g.Move)
	moves[g.MeId] = g.Down

	ns := board.Clone()
	ns.AdvanceBoard(moves)

	newHash := GetZobristKey(tt, ns)

	if newHash == hash1 {
		panic("Snake moved, hashes should not match!")
	}

}
