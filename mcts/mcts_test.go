package mcts

import (
	"fmt"
	"testing"
	"time"

	api "github.com/nosnaws/tiam/battlesnake"
	b "github.com/nosnaws/tiam/board"
)

func TestDecisionHeadToHead(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ e e e
	// _ _ _ _ _ _ _ e e e e
	// _ _ _ _ _ _ _ e _ e _
	// _ _ _ _ _ _ _ h _ _ _
	// _ _ _ _ _ _ h _ _ _ _
	// _ _ _ s s s s f _ _ _
	// _ _ _ s s s f _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []api.Coord{
			{X: 6, Y: 4},
			{X: 6, Y: 3},
			{X: 5, Y: 3},
			{X: 5, Y: 2},
			{X: 4, Y: 2},
			{X: 4, Y: 3},
			{X: 3, Y: 3},
			{X: 3, Y: 2},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{X: 7, Y: 5},
			{X: 7, Y: 6},
			{X: 7, Y: 7},
			{X: 8, Y: 7},
			{X: 8, Y: 8},
			{X: 9, Y: 8},
			{X: 10, Y: 8},
			{X: 10, Y: 7},
			{X: 9, Y: 7},
			{X: 9, Y: 6},
		},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food:   []api.Coord{{X: 7, Y: 3}, {X: 6, Y: 2}},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	move := MCTS(&board, nil)
	if move.Dir == b.Up || move.Dir == b.Right {
		board.Print()
		fmt.Println("selected ", move.Dir)
		panic("Went up or right!")
	}
}

func TestDecision(t *testing.T) {
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
	// f _ _ _ _ _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 10, Y: 5}, {X: 9, Y: 5}, {X: 8, Y: 5}, {X: 7, Y: 5}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 9, Y: 4}, {X: 8, Y: 4}, {X: 7, Y: 4}, {X: 6, Y: 4}},
	}
	three := api.Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []api.Coord{{X: 9, Y: 2}, {X: 8, Y: 2}, {X: 7, Y: 2}},
	}
	four := api.Battlesnake{
		ID:     "four",
		Health: 100,
		Body:   []api.Coord{{X: 5, Y: 6}, {X: 4, Y: 6}, {X: 3, Y: 6}, {X: 2, Y: 6}, {X: 1, Y: 6}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two, three, four},
			Height: 11,
			Width:  11,
			Food:   []api.Coord{{X: 10, Y: 4}, {X: 0, Y: 0}},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	move := MCTS(&board, nil)
	if move.Dir != b.Up {
		board.Print()
		fmt.Println("selected ", move.Dir)
		panic("Did not go up!")
	}
}

func TestDecisionDraw(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _ _ _ _ s _ _
	// _ _ _ _ _ _ _ _ s _ _
	// _ _ _ _ _ _ _ _ s _ _
	// _ _ _ _ _ _ _ _ e _ _
	// _ _ _ _ s s s h f _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ f _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 7, Y: 6}, {X: 6, Y: 6}, {X: 5, Y: 6}, {X: 4, Y: 6}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 8, Y: 7}, {X: 8, Y: 8}, {X: 8, Y: 9}, {X: 8, Y: 10}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food:   []api.Coord{{X: 8, Y: 6}, {X: 1, Y: 2}},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	move := MCTS(&board, nil)
	if move.Dir == b.Right {
		board.Print()
		fmt.Println("selected ", move.Dir)
		panic("Went right for some stupid reason!")
	}
}

func TestDecisionDoesNotSuicide(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _ s s s s s _
	// s s s s s s _ f s s f
	// s _ _ _ s s s s h _ _
	// s _ _ _ s s s e _ _ _
	// s _ _ _ _ _ s s _ _ _
	// s _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []api.Coord{
			{X: 8, Y: 8},
			{X: 8, Y: 9},
			{X: 9, Y: 9},
			{X: 9, Y: 10},
			{X: 8, Y: 10},
			{X: 7, Y: 10},
			{X: 6, Y: 10},
			{X: 5, Y: 10},
			{X: 5, Y: 9},
			{X: 5, Y: 8},
			{X: 6, Y: 8},
			{X: 7, Y: 8},
		},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []api.Coord{
			{X: 7, Y: 7},
			{X: 7, Y: 6},
			{X: 6, Y: 6},
			{X: 6, Y: 7},
			{X: 5, Y: 7},
			{X: 4, Y: 7},
			{X: 4, Y: 8},
			{X: 4, Y: 9},
			{X: 3, Y: 9},
			{X: 2, Y: 9},
			{X: 1, Y: 9},
			{X: 0, Y: 9},
			{X: 0, Y: 8},
			{X: 0, Y: 7},
			{X: 0, Y: 6},
			{X: 0, Y: 5},
		},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food: []api.Coord{
				{X: 7, Y: 9},
				{X: 10, Y: 9},
			},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	move := MCTS(&board, nil)
	if move.Dir == b.Down {
		board.Print()
		fmt.Println("selected ", move.Dir)
		panic("Went down for some stupid reason!")
	}
}

func TestPerformance(t *testing.T) {
	t.Skip()
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ s s s s h _ _ _ _ _
	// _ _ _ _ _ _ _ s s s h
	// _ _ _ _ _ _ s s s h f
	// f _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ s s h _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 10, Y: 5}, {X: 9, Y: 5}, {X: 8, Y: 5}, {X: 7, Y: 5}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 9, Y: 4}, {X: 8, Y: 4}, {X: 7, Y: 4}, {X: 6, Y: 4}},
	}
	three := api.Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []api.Coord{{X: 9, Y: 2}, {X: 8, Y: 2}, {X: 7, Y: 2}},
	}
	four := api.Battlesnake{
		ID:     "four",
		Health: 100,
		Body:   []api.Coord{{X: 5, Y: 6}, {X: 4, Y: 6}, {X: 3, Y: 6}, {X: 2, Y: 6}, {X: 1, Y: 6}},
	}
	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two, three, four},
			Height: 11,
			Width:  11,
			Food:   []api.Coord{{X: 10, Y: 4}, {X: 0, Y: 3}},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	totalRuns := 100
	actualRuns := 0
	maxTime := float64(0)
	totalTime := float64(0)
	for i := 0; i < totalRuns; i++ {
		ns := board.Clone()

		now := time.Now()
		MCTS(&ns, nil)
		after := time.Now()
		actualRuns += 1

		total := float64(after.UnixMilli() - now.UnixMilli())
		totalTime += total
		if maxTime < total {
			maxTime = total
			if maxTime > float64(400) {
				fmt.Println("High duration: ", maxTime)
				break
			}
		}
	}

	fmt.Println("Total runs: ", actualRuns)
	fmt.Println("Average duration: ", float64(totalTime)/float64(totalRuns))
	fmt.Println("Max duration: ", maxTime)
}

func TestPerformanceOpenPosition(t *testing.T) {
	t.Skip()

	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ e _ _ _ _ _ e _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ f _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ _ e _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
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
			Food:   []api.Coord{{X: 5, Y: 5}},
		},
		You: me,
	}
	board := b.BuildBoard(state)

	totalRuns := 500
	actualRuns := 0
	maxTime := float64(0)
	totalTime := float64(0)
	for i := 0; i < totalRuns; i++ {
		ns := board.Clone()

		now := time.Now()
		MCTS(&ns, nil)
		after := time.Now()
		actualRuns += 1

		total := float64(after.UnixMilli() - now.UnixMilli())
		totalTime += total
		if maxTime < total {
			maxTime = total
			if maxTime > float64(450) {
				fmt.Println("High duration: ", maxTime)
				break
			}
		}
	}

	fmt.Println("Total runs: ", actualRuns)
	fmt.Println("Average duration: ", float64(totalTime)/float64(totalRuns))
	fmt.Println("Max duration: ", maxTime)
}
