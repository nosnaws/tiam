package mcts

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
)

func TestStrategicRollout(t *testing.T) {
	//t.Skip()
	// e _ _
	// _ f _
	// _ _ h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   api.Coord{X: 0, Y: 2},
		Body:   []api.Coord{{X: 0, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 2}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 1, Y: 1}},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]
	twoId := board.Ids["two"]

	fmt.Println("Running strategic rollout")
	StrategicRollout(&board)

	fmt.Println(board)
	if board.Healths[id] > 1 && board.Healths[twoId] > 1 {
		fmt.Println(board)
		panic("game did not end!")
	}
}

func TestRandomRollout(t *testing.T) {
	//t.Skip()
	// e _ _
	// _ f _
	// _ _ h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   api.Coord{X: 0, Y: 2},
		Body:   []api.Coord{{X: 0, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 2}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 1, Y: 1}},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]
	twoId := board.Ids["two"]

	fmt.Println("Running random rollout")
	RandomRollout(&board)

	fmt.Println(board)
	if board.Healths[id] > 1 && board.Healths[twoId] > 1 {
		fmt.Println(board)
		panic("game did not end!")
	}
}

func TestRandomRolloutWrapped(t *testing.T) {
	//t.Skip()
	// e _ _
	// _ f _
	// _ _ h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   api.Coord{X: 0, Y: 2},
		Body:   []api.Coord{{X: 0, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 2}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 1, Y: 1}},
		},
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
			},
		},
		You: me,
	}
	board := b.BuildBoard(state)
	id := board.Ids["me"]
	twoId := board.Ids["two"]

	fmt.Println("Running random wrapped rollout")
	RandomRollout(&board)

	fmt.Println(board)
	if board.Healths[id] > 1 && board.Healths[twoId] > 1 {
		fmt.Println(board)
		panic("game did not end!")
	}
}
