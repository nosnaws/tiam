package bitboard

import (
	"math/rand"
	"testing"
	"time"

	api "github.com/nosnaws/tiam/battlesnake"
)

func TestRandomRollout(t *testing.T) {
	// _ _ _ _ _ _ f _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ e _ _ _ _ _ e _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ f _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ _ e _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ f _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 3}, {X: 2, Y: 3}, {X: 2, Y: 3}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 8}, {X: 2, Y: 8}, {X: 2, Y: 8}},
	}
	three := api.Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []api.Coord{{X: 8, Y: 8}, {X: 8, Y: 8}, {X: 8, Y: 8}},
	}
	four := api.Battlesnake{
		ID:     "four",
		Health: 100,
		Body:   []api.Coord{{X: 8, Y: 3}, {X: 8, Y: 3}, {X: 8, Y: 3}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two, three, four},
			Height: 11,
			Width:  11,
			Food:   []api.Coord{{X: 5, Y: 5}, {X: 4, Y: 0}, {X: 6, Y: 10}},
		},
		You: me,
	}

	board := CreateBitBoard(state)

	rand := rand.New(rand.NewSource(time.Now().Unix()))
	board.RandomPlayout(10, rand)

	board.Print()

	ns := board.Clone()
	ns.Print()
}

func BenchmarkRandomRollout(b *testing.B) {
	// _ _ _ _ _ _ f _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ e _ _ _ _ _ e _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ f _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ _ e _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ f _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 3}, {X: 2, Y: 3}, {X: 2, Y: 3}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 8}, {X: 2, Y: 8}, {X: 2, Y: 8}},
	}
	three := api.Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []api.Coord{{X: 8, Y: 8}, {X: 8, Y: 8}, {X: 8, Y: 8}},
	}
	four := api.Battlesnake{
		ID:     "four",
		Health: 100,
		Body:   []api.Coord{{X: 8, Y: 3}, {X: 8, Y: 3}, {X: 8, Y: 3}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two, three, four},
			Height: 11,
			Width:  11,
			Food:   []api.Coord{{X: 5, Y: 5}, {X: 4, Y: 0}, {X: 6, Y: 10}},
		},
		You: me,
	}

	board := CreateBitBoard(state)

	rand := rand.New(rand.NewSource(time.Now().Unix()))
	for n := 0; n < b.N; n++ {
		ns := board.Clone()
		ns.RandomPlayout(10, rand)
	}
}
