package bitboard

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
)

func TestCreateSnake(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ f
	// s s h
	me := api.Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}

	snake := createSnake(me, 3)

	if snake.Length != 3 {
		panic("length is not 3")
	}

	if snake.health != 50 {
		panic("health is not 50")
	}

	if snake.GetHeadIndex() != 2 {
		panic("head not in right spot")
	}

	if snake.getTailIndex() != 0 {
		panic("tail not in right spot")
	}
}

func TestCreateSnakeEaten(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ f
	// s s h
	me := api.Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}},
	}

	snake := createSnake(me, 3)

	if snake.Length != 4 {
		panic("length is not 3")
	}

	if snake.health != 50 {
		panic("health is not 50")
	}

	if snake.GetHeadIndex() != 2 {
		panic("head not in right spot")
	}

	if snake.getTailIndex() != 0 {
		panic("tail not in right spot")
	}
}

func TestCreatesLargeSnake(t *testing.T) {
	//{
	//"latency": "323.376",
	//"id": "f93c4472-3e67-469b-83db-67bceb10c1d8",
	//"health": 91,
	//"length": 8,
	//"shout": "",
	//"head": {
	//"y": 3,
	//"x": 6
	//},
	//"customizations": {
	//"color": "#002080",
	//"tail": "fat-rattle",
	//"head": "evil"
	//},
	//"body": [
	//{
	//"y": 3,
	//"x": 6
	//},
	//{
	//"y": 3,
	//"x": 7
	//},
	//{
	//"y": 3,
	//"x": 8
	//},
	//{
	//"y": 3,
	//"x": 9
	//},
	//{
	//"y": 3,
	//"x": 10
	//},
	//{
	//"y": 4,
	//"x": 10
	//},
	//{
	//"y": 5,
	//"x": 10
	//},
	//{
	//"y": 6,
	//"x": 10
	//}
	//],
	//"name": "main",
	//"squad": ""
	//},

	fmt.Println("BEGIN")
	s := api.Battlesnake{
		ID:     "me",
		Health: 91,
		Length: 8,
		Head:   api.Coord{X: 3, Y: 6},
		Body: []api.Coord{
			{X: 3, Y: 6},
			{X: 3, Y: 7},
			{X: 3, Y: 8},
			{X: 3, Y: 9},
			{X: 3, Y: 10},
			{X: 4, Y: 10},
			{X: 5, Y: 10},
			{X: 6, Y: 10},
		},
	}
	state := api.GameState{
		Turn: 60,
		Board: api.Board{
			Snakes: []api.Battlesnake{s},
			Height: 11,
			Width:  11,
		},
		You: s,
	}

	b := CreateBitBoard(state)
	snake := createSnake(s, 11)

	if snake.Length != 8 {
		panic("length is not 8")
	}

	if snake.health != 91 {
		panic("health is not 91")
	}

	if snake.GetHeadIndex() != 69 {
		snake.print()
		b.Print()
		b.printBoard(snake.board)
		panic("head not in right spot")
	}

	if snake.getTailIndex() != 116 {
		snake.print()
		b.Print()
		b.printBoard(snake.board)
		panic("tail not in right spot")
	}

}
