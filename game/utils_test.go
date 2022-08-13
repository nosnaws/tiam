package game

import (
	"fmt"
	"testing"
)

func TestIndexInDirection(t *testing.T) {

	if IndexInDirection(Right, 0, 3, 3, false) != 1 {
		panic("Right should be TileIndex 1")
	}

	if IndexInDirection(Up, 0, 3, 3, false) != 3 {
		panic("Up should be TileIndex 3")
	}

	if IndexInDirection(Down, 4, 3, 3, false) != 1 {
		panic("Down should be TileIndex 1")
	}

	if IndexInDirection(Left, 1, 3, 3, false) != 0 {
		fmt.Println(IndexInDirection(Left, 1, 3, 3, false))
		panic("Left should be TileIndex 0")
	}

}

func TestMoveToPoint(t *testing.T) {
	p := Point{X: -1, Y: 0}
	if moveToPoint(Left) != p {
		panic("Not left!")
	}
	p = Point{X: 1, Y: 0}
	if moveToPoint(Right) != p {
		panic("Not right!")
	}
	p = Point{X: 0, Y: 1}
	if moveToPoint(Up) != p {
		panic("Not up!")
	}
	p = Point{X: 0, Y: -1}
	if moveToPoint(Down) != p {
		panic("Not down!")
	}
}

func TestRandomCartesianProductWrapped(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// _ 2 _ 3 _
	// _ _ _ _ _
	// _ 4 _ h _
	// _ _ _ _ _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 3, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 1}},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 3}, {X: 1, Y: 3}, {X: 1, Y: 3}},
	}
	three := Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []Coord{{X: 3, Y: 3}, {X: 3, Y: 3}, {X: 3, Y: 3}},
	}
	four := Battlesnake{
		ID:     "four",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}},
	}
	state := GameState{
		Turn: 1,
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
		Board: Board{
			Snakes: []Battlesnake{me, two, three, four},
			Height: 5,
			Width:  5,
		},

		You: me,
	}
	board := BuildBoard(state)

	cartMoves := GetCartesianProductOfMoves(board)

	if len(cartMoves) != 256 {
		fmt.Println("total moves ", len(cartMoves))
		panic("Did not create 256 possible states 4 to the 4nd")
	}
}

func TestRandomCartesianProduct(t *testing.T) {
	//t.Skip()
	// e _ _
	// _ f _
	// _ _ h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   Coord{X: 0, Y: 2},
		Body:   []Coord{{X: 0, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 2}},
	}
	state := GameState{
		Turn: 1,
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 1, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)

	cartMoves := GetCartesianProductOfMoves(board)

	if len(cartMoves) != 4 {
		fmt.Println("total moves ", len(cartMoves))
		panic("Did not create 4 possible states 2 to the 2nd")
	}
}
