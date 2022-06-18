package main

import (
	"fmt"
	"testing"
)

//func TestVoronoi(t *testing.T) {
//me := Battlesnake{
//// Length 3, facing right
//ID:   "me",
//Head: Coord{X: 0, Y: 0},
//Body: []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
//}
//s2 := Battlesnake{
//// Length 3, facing right
//ID:   "s2",
//Head: Coord{X: 3, Y: 4},
//Body: []Coord{{X: 3, Y: 4}, {X: 3, Y: 3}, {X: 3, Y: 2}},
//}
//state := GameState{
//Board: Board{
//Snakes: []Battlesnake{me, s2},
//Height: 5,
//Width:  5,
//},
//You: me,
//}

//game := BuildBoard(state)
//v := voronoi(&game, false, 5)
//fmt.Println(v)

//if v[coordToPoint(me.Head)] != 4 {
//panic(fmt.Sprintf("voronoi %d != 8", v[coordToPoint(me.Head)]))
//}
//if v[coordToPoint(s2.Head)] != 6 {
//panic(fmt.Sprintf("voronoi %d != 11", v[coordToPoint(s2.Head)]))
//}
//}

func TestGetEdges(t *testing.T) {

	me := Battlesnake{
		// Length 3, facing right
		ID:   "me",
		Head: Coord{X: 0, Y: 0},
		Body: []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	s2 := Battlesnake{
		// Length 3, facing right
		ID:   "s2",
		Head: Coord{X: 0, Y: 1},
		Body: []Coord{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me, s2},
			Height: 5,
			Width:  5,
		},
		You: me,
	}

	game := BuildBoard(state)
	edges := GetEdges(coordToPoint(me.Head), &game, true)

	fmt.Println(edges)
	if len(edges) != 2 {
		panic("did not find correct edges")
	}
}
