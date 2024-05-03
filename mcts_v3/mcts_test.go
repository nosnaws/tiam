package mctsv3

import (
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
)

//func Test(t *testing.T) {
////t.Skip()
//// _ _ _
//// _ _ _
//// s s h
//me := api.Battlesnake{
//// Length 3, facing right
//ID:     "me",
//Health: 100,
//Head:   api.Coord{X: 2, Y: 0},
//Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
//}
//state := api.GameState{
//Turn: 3,
//Board: api.Board{
//Snakes: []api.Battlesnake{me},
//Height: 3,
//Width:  3,
//},
//You: me,
//}

//board := bitboard.CreateBitBoard(state)

//node := createNode(&board, []bitboard.SnakeMove{}, nil)
//expand(node, &board)
//if len(node.children) != 1 {
//panic("wrong number of children!")
//}

//if len(node.children[0].reward[bitboard.MeId]) != 2 {
//panic("wrong number of rewards")
//}

//if node.children[0].prevMoves[0].Dir != bitboard.Up {
//panic("incorrect previous moves")
//}

//}

func TestMCTS(t *testing.T) {
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ e _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ e _ _ _ _ e _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _

	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 8},
		Body: []api.Coord{
			{X: 2, Y: 8},
			{X: 2, Y: 8},
			{X: 2, Y: 8},
		},
	}
	s1 := api.Battlesnake{
		ID:     "s1",
		Health: 100,
		Head:   api.Coord{X: 7, Y: 8},
		Body: []api.Coord{
			{X: 7, Y: 8},
			{X: 7, Y: 8},
			{X: 7, Y: 8},
		},
	}
	s2 := api.Battlesnake{
		ID:     "s2",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 3},
		Body: []api.Coord{
			{X: 2, Y: 3},
			{X: 2, Y: 3},
			{X: 2, Y: 3},
		},
	}
	s3 := api.Battlesnake{
		ID:     "s3",
		Health: 100,
		Head:   api.Coord{X: 7, Y: 3},
		Body: []api.Coord{
			{X: 7, Y: 3},
			{X: 7, Y: 3},
			{X: 7, Y: 3},
		},
	}

	state := api.GameState{
		Turn: 1,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, s1, s2, s3},
			Height: 11,
			Width:  11,
		},
		You: me,
	}

	node := ChooseNewRoot(nil, state)
	MCTS(node, 10, 5, 20)
}

func BenchmarkMCTS(b *testing.B) {
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ e _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ e _ _ _ _ e _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _

	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 8},
		Body: []api.Coord{
			{X: 2, Y: 8},
			{X: 2, Y: 8},
			{X: 2, Y: 8},
		},
	}
	s1 := api.Battlesnake{
		ID:     "s1",
		Health: 100,
		Head:   api.Coord{X: 7, Y: 8},
		Body: []api.Coord{
			{X: 7, Y: 8},
			{X: 7, Y: 8},
			{X: 7, Y: 8},
		},
	}
	s2 := api.Battlesnake{
		ID:     "s2",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 3},
		Body: []api.Coord{
			{X: 2, Y: 3},
			{X: 2, Y: 3},
			{X: 2, Y: 3},
		},
	}
	s3 := api.Battlesnake{
		ID:     "s3",
		Health: 100,
		Head:   api.Coord{X: 7, Y: 3},
		Body: []api.Coord{
			{X: 7, Y: 3},
			{X: 7, Y: 3},
			{X: 7, Y: 3},
		},
	}

	state := api.GameState{
		Turn: 1,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, s1, s2, s3},
			Height: 11,
			Width:  11,
		},
		You: me,
	}

	for n := 0; n < b.N; n++ {
		node := ChooseNewRoot(nil, state)
		MCTS(node, 10, 5, 20)
	}
}
