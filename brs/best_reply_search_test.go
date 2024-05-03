package brain

import (
	"fmt"
	"math"
	"testing"

	"github.com/nosnaws/tiam/game"
)

func TestBRSDoesntGoBackward(t *testing.T) {
	//t.Skip()
	// _ _ _ _ f _ f _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// s _ _ _ _ f _ _ _ _ _
	// s s s _ _ _ _ _ _ _ _
	// _ e s _ _ _ _ _ _ _ _
	// _ f h _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	me := game.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []game.Coord{
			{2, 2},
			{2, 3},
			{2, 4},
		},
	}
	two := game.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []game.Coord{
			{1, 3},
			{1, 4},
			{0, 4},
			{0, 5},
		},
	}
	state := game.GameState{
		Turn: 0,
		Board: game.Board{
			Snakes: []game.Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food: []game.Coord{
				{6, 10},
				{5, 5},
				{1, 2},
				{4, 10},
			},
		},
		You: me,
	}
	board := game.BuildBoard(state)
	brs := CreateBRSGame(&board)

	move := brs.BRS(&board, game.Left, 9, math.Inf(-1), math.Inf(1), true)
	if move.Move == game.Left {
		board.Print()
		fmt.Println("selected ", move)
		panic("Went up or right!")
	}
}

func TestBRSTrap(t *testing.T) {
	//t.Skip()
	// s s s _ _ _ _ _ _ _ _
	// s _ s _ _ _ _ _ _ _ _
	// s h s _ _ _ f _ _ _ _
	// _ _ s s _ f _ _ _ _ _
	// _ _ _ s _ _ _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ e s s s s _ _ _
	// _ _ _ _ _ _ _ s _ _ _
	// f f _ _ _ _ _ s _ _ _
	// _ _ f _ _ _ s s _ _ _
	// _ _ _ _ _ _ s s s _ _
	me := game.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []game.Coord{
			{1, 8},
			{0, 8},
			{0, 9},
			{0, 10},
			{1, 10},
			{2, 10},
			{2, 9},
			{2, 8},
			{2, 7},
			{3, 7},
			{3, 6},
		},
	}
	two := game.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []game.Coord{
			{3, 4},
			{4, 4},
			{5, 4},
			{6, 4},
			{7, 4},
			{7, 3},
			{7, 2},
			{7, 1},
			{6, 1},
			{6, 0},
			{7, 0},
			{8, 0},
		},
	}
	state := game.GameState{
		Turn: 0,
		Board: game.Board{
			Snakes: []game.Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food: []game.Coord{
				{5, 8},
				{4, 7},
				{0, 2},
				{1, 2},
				{2, 1},
			},
		},
		You: me,
	}
	board := game.BuildBoard(state)
	brs := BRSGame{}

	fmt.Println("BEGIN $$$$$$$$$$$$$$$$$$$$$$$$")
	move := brs.BRS(&board, game.Left, 3, math.Inf(-1), math.Inf(1), true)
	fmt.Println("SELECTED", move)
	if move.Move != game.Down {
		board.Print()
		fmt.Println("selected ", move)
		panic("Went up or right!")
	}
}

func TestBRSTrapSmall(t *testing.T) {
	//t.Skip()
	// s s s _ _
	// s _ s _ _
	// s h s s x
	// _ _ _ e x
	// _ _ _ x x
	me := game.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []game.Coord{
			{1, 2},
			{0, 2},
			{0, 3},
			{0, 4},
			{1, 4},
			{2, 4},
			{2, 3},
			{2, 2},
			{2, 2},
			{3, 2},
		},
	}
	two := game.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []game.Coord{
			{3, 1},
			{3, 0},
			{4, 0},
			{4, 1},
			{4, 2},
		},
	}
	state := game.GameState{
		Turn: 0,
		Board: game.Board{
			Snakes: []game.Battlesnake{me, two},
			Height: 11,
			Width:  11,
		},
		You: me,
	}
	board := game.BuildBoard(state)

	brs := BRSGame{}

	move := brs.BRS(&board, game.Left, 9, math.Inf(-1), math.Inf(1), true)
	if move.Move != game.Down {
		board.Print()
		fmt.Println("selected ", move)
		panic("Went up or right!")
	}

}
