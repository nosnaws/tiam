package game

import (
	"fmt"
	"testing"
)

func TestBoardCreation(t *testing.T) {
	//t.Skip()
	// Arrange
	me := Battlesnake{
		// Length 3, facing right
		Head: Coord{X: 2, Y: 0},
		Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes:  []Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []Coord{{X: 2, Y: 0}},
			Food:    []Coord{{X: 2, Y: 1}},
		},
		You: me,
	}

	board := BuildBoard(state)

	if board.list[2].id != 1 {
		panic("YouId is not 1")
	}

	if board.list[2].IsHazard() != true {
		panic("Did not create hazard")
	}

	if board.list[5].IsFood() != true {
		panic("Did not create food")
	}

	head := board.list[2]
	tail := board.list[head.GetIdx()]
	current := tail
	for current != head {
		current = board.list[current.GetIdx()]
	}

	if current != head {
		panic("Snake does not loop to head!")
	}
}

func TestKill(t *testing.T) {
	//t.Skip()
	// Arrange
	me := Battlesnake{
		// Length 3, facing right
		ID:   "me",
		Head: Coord{X: 2, Y: 0},
		Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes:  []Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []Coord{{X: 2, Y: 0}},
			Food:    []Coord{{X: 2, Y: 1}},
		},
		You: me,
	}

	board := BuildBoard(state)
	id := board.ids["me"]
	board.kill(id)

	if board.list[2].id == 1 {
		panic("Snake was not removed!")
	}

	if board.heads[id] != 0 {
		panic("Snake head should be set to 0 in head map!")
	}
	if board.lengths[id] != 0 {
		panic("Snake length should be set to 0 in length map!")
	}
	if board.healths[id] != 0 {
		panic("Snake health should be set to 0 in health map!")
	}
}

func TestAdvanceBoardMoving(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Up}}
	board.advanceBoard(moves)

	if board.list[5].IsSnakeHead() != true {
		fmt.Println(board)
		fmt.Println(board.list[5].idx)
		panic("Did not move snake up!")
	}

	moves = []snakeMove{{id: id, dir: Left}}
	board.advanceBoard(moves)

	if board.list[4].IsSnakeHead() != true {
		fmt.Println(board)
		panic("Did not move snake left!")
	}

	moves = []snakeMove{{id: id, dir: Down}}
	board.advanceBoard(moves)

	if board.list[1].IsSnakeHead() != true {
		panic("Did not move snake down!")
	}

	moves = []snakeMove{{id: id, dir: Right}}
	board.advanceBoard(moves)

	if board.list[2].IsSnakeHead() != true {
		panic("Did not move snake right!")
	}
}

func TestAdvanceBoardTurnDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// s s h
	me := Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}

	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Up}}
	board.advanceBoard(moves)

	if board.healths[id] != 99 {
		fmt.Println(board.healths[id])
		panic("Did not decrement health properly!")
	}
}

func TestAdvanceBoardHazardDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ f
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes:  []Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Up}}
	board.advanceBoard(moves)

	if board.healths[id] > 0 {
		panic("Snake did not die!")
	}
}

func TestAdvanceBoardEatFood(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ f
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Up}}
	board.advanceBoard(moves)

	if board.healths[id] != 100 {
		panic("Snake did not eat!")
	}
}

func TestAdvanceBoardOutOfBounds(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Right}}
	board.advanceBoard(moves)

	if board.list[2].IsSnakeSegment() != false && board.list[1].IsSnakeSegment() != false && board.list[0].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}
}

func TestAdvanceBoardHeadCollision(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ f
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	enemyId := board.ids["enemy"]

	moves := []snakeMove{{id: id, dir: Up}, {id: enemyId, dir: Down}}
	board.advanceBoard(moves)

	if board.healths[id] != 0 {
		panic("Did not remove me from board!")
	}

	if board.healths[enemyId] != 100 {
		fmt.Println(board.healths)
		panic("Should not have removed enemy!")
	}
}

func TestAdvanceBoardSnakeCollision(t *testing.T) {
	//t.Skip()
	// o o _
	// o e _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	enemyId := board.ids["enemy"]

	moves := []snakeMove{{id: id, dir: Up}, {id: enemyId, dir: Down}}
	board.advanceBoard(moves)

	if board.healths[enemyId] != 0 {
		fmt.Println(board)
		fmt.Println(board.healths[enemyId])
		panic("Should have removed enemy from the board!")
	}

	if board.healths[id] != 99 {
		panic("Should not have removed me!")
	}
}

func TestAdvanceBoardSelfCollision(t *testing.T) {
	//t.Skip()
	// _ s _
	// _ s h
	// _ s s
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Left}}
	board.advanceBoard(moves)

	if board.healths[id] != 0 {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardFollowTail(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ s s
	// _ s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Up}}
	board.advanceBoard(moves)

	if board.healths[id] < 1 {
		panic("Snake was killed!")
	}

	if board.list[5].IsSnakeHead() != true {
		panic("Did not move snake!")
	}

	// Follow other snake tail
	// _ e e
	// _ s e
	// _ s h
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.ids["me"]
	enemyId := board.ids["enemy"]

	moves = []snakeMove{{id: id, dir: Up}, {id: enemyId, dir: Left}}
	fmt.Println("advancing board")
	board.advanceBoard(moves)
	fmt.Println("board advanced")

	if board.healths[id] < 1 {
		panic("Snake was killed!")
	}

	if board.list[5].IsSnakeHead() != true {
		fmt.Println(board)
		panic("Did not move snake!")
	}

	// Follow other snake tail
	// f e e
	// _ s e
	// _ s h
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	}
	enemy = Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 0, Y: 2}},
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.ids["me"]
	enemyId = board.ids["enemy"]

	moves = []snakeMove{{id: id, dir: Up}, {id: enemyId, dir: Left}}
	fmt.Println("advancing board")
	board.advanceBoard(moves)
	fmt.Println("board advanced")

	if board.healths[id] != 0 {
		panic("Snake was not killed!")
	}

	if board.list[5].IsSnakeHead() != false {
		fmt.Println(board)
		panic("Did not remove snake!")
	}
}

func TestAdvanceBoardMoveOnNeck(t *testing.T) {
	//t.Skip()
	// _ s _
	// _ s h
	// _ s s
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Down}}
	board.advanceBoard(moves)

	if board.healths[id] != 0 {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardWrapped(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ h
	// _ s s
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	moves := []snakeMove{{id: id, dir: Right}}
	board.advanceBoard(moves)

	if board.healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.list[3].IsSnakeHead() != true {
		panic("Did not wrapped to right!")
	}

	// _ _ h
	// _ _ s
	// _ _ s
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.ids["me"]

	moves = []snakeMove{{id: id, dir: Up}}
	board.advanceBoard(moves)

	if board.healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.list[2].IsSnakeHead() != true {
		panic("Did not wrapped to bottom!")
	}

	// _ s _
	// _ s _
	// _ h _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.ids["me"]

	moves = []snakeMove{{id: id, dir: Down}}
	board.advanceBoard(moves)

	if board.healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.list[7].IsSnakeHead() != true {
		panic("Did not wrapped to top!")
	}

	// _ _ _
	// s s h
	// _ _ _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 1}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.ids["me"]

	moves = []snakeMove{{id: id, dir: Right}}
	board.advanceBoard(moves)

	if board.healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.list[3].IsSnakeHead() != true {
		panic("Did not wrapped to left!")
	}
}

func TestAdvanceBoardCrazyStuff(t *testing.T) {
	//t.Skip()
	//   s . . f . . . . s 3 s
	//   . . . . . . . . s . .
	//   . . . . . . . . . . .
	//   . f . . . . . . . . .
	//   . . . . . . . . . . .
	//   . . f . . s . . . . .
	//   . . . . 2 s . . . . .
	//   . . . . . . . . . . .
	//   . . . . . . . . . . .
	//   . . . . . . . . . 0 .
	//   . . . . . . . . 1 s s
	zero := Battlesnake{
		ID:     "0",
		Health: 100,
		Body:   []Coord{{X: 9, Y: 1}, {X: 9, Y: 0}, {X: 10, Y: 0}},
	}
	one := Battlesnake{
		ID:     "1",
		Health: 100,
		Body:   []Coord{{X: 8, Y: 0}, {X: 8, Y: 10}, {X: 8, Y: 9}},
	}
	two := Battlesnake{
		ID:     "2",
		Health: 100,
		Body:   []Coord{{X: 4, Y: 4}, {X: 5, Y: 4}, {X: 5, Y: 5}},
	}
	three := Battlesnake{
		ID:     "3",
		Health: 100,
		Body:   []Coord{{X: 9, Y: 10}, {X: 10, Y: 10}, {X: 0, Y: 10}},
	}

	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{zero, one, two, three},
			Height: 11,
			Width:  11,
		},
		You: zero,
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := BuildBoard(state)
	id0 := board.ids["0"]
	id1 := board.ids["1"]
	id2 := board.ids["2"]
	id3 := board.ids["3"]

	moves := []snakeMove{{id: id0, dir: Up}, {id: id1, dir: Right}, {id: id2, dir: Left}, {id: id3, dir: Up}}
	board.advanceBoard(moves)

	if board.healths[id1] > 0 {
		panic("Should not have killed id1!")
	}

	if board.healths[id3] > 0 {
		panic("Should not have killed id3!")
	}

	snakeId, ok := board.list[9].GetSnakeId()
	if board.list[9].IsSnakeSegment() != true || !ok || snakeId != id0 {
		panic("lost snake id0 neck!")
	}
}
