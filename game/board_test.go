package game

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRandomRollout(t *testing.T) {
	t.Skip()
	// e s s
	// _ f _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   Coord{X: 0, Y: 2},
		Body:   []Coord{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 1, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	twoId := board.ids["two"]

	fmt.Println("Running random rollout")
	board.RandomRollout()

	fmt.Println(board)
	if board.Healths[id] > 1 && board.Healths[twoId] > 1 {
		fmt.Println(board)
		panic("game did not end!")
	}
}

func TestRandomRolloutWrapped(t *testing.T) {
	t.Skip()
	// e s s
	// _ f _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   Coord{X: 0, Y: 2},
		Body:   []Coord{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 1, Y: 1}},
		},
		Game: Game{
			Ruleset: Ruleset{
				Name: "wrapped",
			},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	twoId := board.ids["two"]

	fmt.Println("Running random wrapped rollout")
	board.RandomRollout()

	fmt.Println(board)
	if board.Healths[id] > 1 && board.Healths[twoId] > 1 {
		fmt.Println(board)
		panic("game did not end!")
	}
}

func TestGetSnakeMoves(t *testing.T) {
	t.Skip()
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

	moves := board.GetMovesForSnake(id)

	if len(moves) != 1 {
		fmt.Println(moves)
		panic("Should only be able to move up!")
	}

	// wrapped!
	// _ _ _
	// _ _ _
	// s s h
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
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

	moves = board.GetMovesForSnake(id)

	if len(moves) != 3 {
		fmt.Println(moves)
		panic("Should be able to go up,down, and right!")
	}

	// wrapped with snake eating
	// f e s
	// s _ s
	// s h _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 1, Y: 0},
		Body:   []Coord{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   Coord{X: 1, Y: 2},
		Body:   []Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
	}
	state = GameState{
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []Coord{{X: 0, Y: 2}},
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
	enemy := board.ids["two"]

	board.AdvanceBoard([]SnakeMove{{Id: id, Dir: Right}, {Id: enemy, Dir: Left}})
	fmt.Println(board)
	moves = board.GetMovesForSnake(id)

	if len(moves) != 2 {
		fmt.Println(moves)
		panic("Should be able to go right and up!")
	}
}

func TestBoardCreation(t *testing.T) {
	t.Skip()
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
	t.Skip()
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

	if board.Heads[id] != 0 {
		panic("Snake head should be set to 0 in head map!")
	}
	if board.Lengths[id] != 0 {
		panic("Snake length should be set to 0 in length map!")
	}
	if board.Healths[id] != 0 {
		panic("Snake health should be set to 0 in health map!")
	}
}

func TestAdvanceBoardMoving(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.list[5].IsSnakeHead() != true {
		fmt.Println(board)
		fmt.Println(board.list[5].idx)
		panic("Did not move snake up!")
	}

	moves = []SnakeMove{{Id: id, Dir: Left}}
	board.AdvanceBoard(moves)

	if board.list[4].IsSnakeHead() != true {
		fmt.Println(board)
		panic("Did not move snake left!")
	}

	moves = []SnakeMove{{Id: id, Dir: Down}}
	board.AdvanceBoard(moves)

	if board.list[1].IsSnakeHead() != true {
		panic("Did not move snake down!")
	}

	moves = []SnakeMove{{Id: id, Dir: Right}}
	board.AdvanceBoard(moves)

	if board.list[2].IsSnakeHead() != true {
		panic("Did not move snake right!")
	}
}

func TestAdvanceBoardTurnDamage(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		fmt.Println(board.Healths[id])
		panic("Did not decrement health properly!")
	}
}

func TestAdvanceBoardHazardDamage(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.Healths[id] > 0 {
		panic("Snake did not die!")
	}
}

func TestAdvanceBoardEatFood(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 100 {
		panic("Snake did not eat!")
	}
}

func TestAdvanceBoardOutOfBounds(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Right}}
	board.AdvanceBoard(moves)

	if board.list[2].IsSnakeSegment() != false && board.list[1].IsSnakeSegment() != false && board.list[0].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}
}

func TestAdvanceBoardHeadCollision(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}, {Id: enemyId, Dir: Down}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 0 {
		panic("Did not remove me from board!")
	}

	if board.Healths[enemyId] != 100 {
		fmt.Println(board.Healths)
		panic("Should not have removed enemy!")
	}
}

func TestAdvanceBoardSnakeCollision(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}, {Id: enemyId, Dir: Down}}
	board.AdvanceBoard(moves)

	if board.Healths[enemyId] != 0 {
		fmt.Println(board)
		fmt.Println(board.Healths[enemyId])
		panic("Should have removed enemy from the board!")
	}

	if board.Healths[id] != 99 {
		panic("Should not have removed me!")
	}
}

func TestAdvanceBoardSelfCollision(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Left}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 0 {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardFollowTail(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.Healths[id] < 1 {
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

	moves = []SnakeMove{{Id: id, Dir: Up}, {Id: enemyId, Dir: Left}}
	fmt.Println("advancing board")
	board.AdvanceBoard(moves)
	fmt.Println("board advanced")

	if board.Healths[id] < 1 {
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

	moves = []SnakeMove{{Id: id, Dir: Up}, {Id: enemyId, Dir: Left}}
	fmt.Println("advancing board")
	board.AdvanceBoard(moves)
	fmt.Println("board advanced")

	if board.Healths[id] != 0 {
		panic("Snake was not killed!")
	}

	if board.list[5].IsSnakeHead() != false {
		fmt.Println(board)
		panic("Did not remove snake!")
	}
}

func TestAdvanceBoardMoveOnNeck(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Down}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 0 {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardWrapped(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id, Dir: Right}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
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

	moves = []SnakeMove{{Id: id, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
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

	moves = []SnakeMove{{Id: id, Dir: Down}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
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

	moves = []SnakeMove{{Id: id, Dir: Right}}
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.list[3].IsSnakeHead() != true {
		panic("Did not wrapped to left!")
	}
}

func TestAdvanceBoardCrazyStuff(t *testing.T) {
	t.Skip()
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

	moves := []SnakeMove{{Id: id0, Dir: Up}, {Id: id1, Dir: Right}, {Id: id2, Dir: Left}, {Id: id3, Dir: Up}}
	board.AdvanceBoard(moves)

	if board.Healths[id1] > 0 {
		panic("Should not have killed id1!")
	}

	if board.Healths[id3] > 0 {
		panic("Should not have killed id3!")
	}

	snakeId, ok := board.list[9].GetSnakeId()
	if board.list[9].IsSnakeSegment() != true || !ok || snakeId != id0 {
		panic("lost snake id0 neck!")
	}
}

func TestBuildBigBoard(t *testing.T) {
	g := []byte("{\"game\":{\"id\":\"c20df634-4097-471b-a558-c6e96ac56620\",\"ruleset\":{\"name\":\"wrapped\",\"version\":\"cli\",\"settings\":{\"foodSpawnChance\":15,\"minimumFood\":1,\"hazardDamagePerTurn\":100,\"hazardMap\":\"\",\"hazardMapAuthor\":\"\",\"royale\":{\"shrinkEveryNTurns\":25},\"squad\":{\"allowBodyCollisions\":false,\"sharedElimination\":false,\"sharedHealth\":false,\"sharedLength\":false}}},\"map\":\"arcade_maze\",\"timeout\":500,\"source\":\"\"},\"turn\":7,\"board\":{\"height\":21,\"width\":19,\"snakes\":[{\"id\":\"9eecc83c-283a-446c-98e0-b1f5e36d9c83\",\"name\":\"tiam\",\"latency\":\"0\",\"health\":93,\"body\":[{\"x\":10,\"y\":10},{\"x\":10,\"y\":9},{\"x\":11,\"y\":9}],\"head\":{\"x\":10,\"y\":10},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}},{\"id\":\"9cca4f6e-8872-42f1-abaf-685887e0697d\",\"name\":\"local\",\"latency\":\"0\",\"health\":93,\"body\":[{\"x\":4,\"y\":14},{\"x\":4,\"y\":13},{\"x\":4,\"y\":12}],\"head\":{\"x\":4,\"y\":14},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}],\"food\":[{\"x\":9,\"y\":11},{\"x\":4,\"y\":17}],\"hazards\":[{\"x\":0,\"y\":20},{\"x\":2,\"y\":20},{\"x\":3,\"y\":20},{\"x\":4,\"y\":20},{\"x\":5,\"y\":20},{\"x\":6,\"y\":20},{\"x\":7,\"y\":20},{\"x\":8,\"y\":20},{\"x\":9,\"y\":20},{\"x\":10,\"y\":20},{\"x\":11,\"y\":20},{\"x\":12,\"y\":20},{\"x\":13,\"y\":20},{\"x\":14,\"y\":20},{\"x\":15,\"y\":20},{\"x\":16,\"y\":20},{\"x\":18,\"y\":20},{\"x\":0,\"y\":19},{\"x\":9,\"y\":19},{\"x\":18,\"y\":19},{\"x\":0,\"y\":18},{\"x\":2,\"y\":18},{\"x\":3,\"y\":18},{\"x\":5,\"y\":18},{\"x\":6,\"y\":18},{\"x\":7,\"y\":18},{\"x\":9,\"y\":18},{\"x\":11,\"y\":18},{\"x\":12,\"y\":18},{\"x\":13,\"y\":18},{\"x\":15,\"y\":18},{\"x\":16,\"y\":18},{\"x\":18,\"y\":18},{\"x\":0,\"y\":17},{\"x\":18,\"y\":17},{\"x\":0,\"y\":16},{\"x\":2,\"y\":16},{\"x\":3,\"y\":16},{\"x\":5,\"y\":16},{\"x\":7,\"y\":16},{\"x\":8,\"y\":16},{\"x\":9,\"y\":16},{\"x\":10,\"y\":16},{\"x\":11,\"y\":16},{\"x\":13,\"y\":16},{\"x\":15,\"y\":16},{\"x\":16,\"y\":16},{\"x\":18,\"y\":16},{\"x\":0,\"y\":15},{\"x\":5,\"y\":15},{\"x\":9,\"y\":15},{\"x\":13,\"y\":15},{\"x\":18,\"y\":15},{\"x\":0,\"y\":14},{\"x\":3,\"y\":14},{\"x\":5,\"y\":14},{\"x\":6,\"y\":14},{\"x\":7,\"y\":14},{\"x\":9,\"y\":14},{\"x\":11,\"y\":14},{\"x\":12,\"y\":14},{\"x\":13,\"y\":14},{\"x\":15,\"y\":14},{\"x\":18,\"y\":14},{\"x\":0,\"y\":13},{\"x\":3,\"y\":13},{\"x\":5,\"y\":13},{\"x\":13,\"y\":13},{\"x\":15,\"y\":13},{\"x\":18,\"y\":13},{\"x\":0,\"y\":12},{\"x\":1,\"y\":12},{\"x\":2,\"y\":12},{\"x\":3,\"y\":12},{\"x\":5,\"y\":12},{\"x\":7,\"y\":12},{\"x\":9,\"y\":12},{\"x\":11,\"y\":12},{\"x\":13,\"y\":12},{\"x\":15,\"y\":12},{\"x\":16,\"y\":12},{\"x\":17,\"y\":12},{\"x\":18,\"y\":12},{\"x\":7,\"y\":11},{\"x\":11,\"y\":11},{\"x\":0,\"y\":10},{\"x\":1,\"y\":10},{\"x\":2,\"y\":10},{\"x\":3,\"y\":10},{\"x\":5,\"y\":10},{\"x\":7,\"y\":10},{\"x\":9,\"y\":10},{\"x\":11,\"y\":10},{\"x\":13,\"y\":10},{\"x\":15,\"y\":10},{\"x\":16,\"y\":10},{\"x\":17,\"y\":10},{\"x\":18,\"y\":10},{\"x\":0,\"y\":9},{\"x\":3,\"y\":9},{\"x\":5,\"y\":9},{\"x\":13,\"y\":9},{\"x\":15,\"y\":9},{\"x\":18,\"y\":9},{\"x\":0,\"y\":8},{\"x\":3,\"y\":8},{\"x\":5,\"y\":8},{\"x\":7,\"y\":8},{\"x\":8,\"y\":8},{\"x\":9,\"y\":8},{\"x\":10,\"y\":8},{\"x\":11,\"y\":8},{\"x\":13,\"y\":8},{\"x\":15,\"y\":8},{\"x\":18,\"y\":8},{\"x\":0,\"y\":7},{\"x\":9,\"y\":7},{\"x\":18,\"y\":7},{\"x\":0,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":7,\"y\":6},{\"x\":9,\"y\":6},{\"x\":11,\"y\":6},{\"x\":12,\"y\":6},{\"x\":13,\"y\":6},{\"x\":15,\"y\":6},{\"x\":16,\"y\":6},{\"x\":18,\"y\":6},{\"x\":0,\"y\":5},{\"x\":3,\"y\":5},{\"x\":15,\"y\":5},{\"x\":18,\"y\":5},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":3,\"y\":4},{\"x\":5,\"y\":4},{\"x\":7,\"y\":4},{\"x\":8,\"y\":4},{\"x\":9,\"y\":4},{\"x\":10,\"y\":4},{\"x\":11,\"y\":4},{\"x\":13,\"y\":4},{\"x\":15,\"y\":4},{\"x\":17,\"y\":4},{\"x\":18,\"y\":4},{\"x\":0,\"y\":3},{\"x\":5,\"y\":3},{\"x\":9,\"y\":3},{\"x\":13,\"y\":3},{\"x\":18,\"y\":3},{\"x\":0,\"y\":2},{\"x\":2,\"y\":2},{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":6,\"y\":2},{\"x\":7,\"y\":2},{\"x\":9,\"y\":2},{\"x\":11,\"y\":2},{\"x\":12,\"y\":2},{\"x\":13,\"y\":2},{\"x\":14,\"y\":2},{\"x\":15,\"y\":2},{\"x\":16,\"y\":2},{\"x\":18,\"y\":2},{\"x\":0,\"y\":1},{\"x\":18,\"y\":1},{\"x\":0,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":6,\"y\":0},{\"x\":7,\"y\":0},{\"x\":8,\"y\":0},{\"x\":9,\"y\":0},{\"x\":10,\"y\":0},{\"x\":11,\"y\":0},{\"x\":12,\"y\":0},{\"x\":13,\"y\":0},{\"x\":14,\"y\":0},{\"x\":15,\"y\":0},{\"x\":16,\"y\":0},{\"x\":18,\"y\":0}]},\"you\":{\"id\":\"9cca4f6e-8872-42f1-abaf-685887e0697d\",\"name\":\"local\",\"latency\":\"0\",\"health\":93,\"body\":[{\"x\":4,\"y\":14},{\"x\":4,\"y\":13},{\"x\":4,\"y\":12}],\"head\":{\"x\":4,\"y\":14},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}}")

	state := GameState{}
	_ = json.Unmarshal(g, &state)

	b := BuildBoard(state)

	if _, ok := b.Heads[MeId]; !ok {
		panic("Did not add me to board!")
	}
}

func TestBuildBigBoard2(t *testing.T) {
	g := []byte("{\"game\":{\"id\":\"676d530f-2d7f-4590-8fda-d15579e5b6fe\",\"ruleset\":{\"name\":\"wrapped\",\"version\":\"cli\",\"settings\":{\"foodSpawnChance\":15,\"minimumFood\":1,\"hazardDamagePerTurn\":100,\"hazardMap\":\"\",\"hazardMapAuthor\":\"\",\"royale\":{\"shrinkEveryNTurns\":25},\"squad\":{\"allowBodyCollisions\":false,\"sharedElimination\":false,\"sharedHealth\":false,\"sharedLength\":false}}},\"map\":\"arcade_maze\",\"timeout\":500,\"source\":\"\"},\"turn\":12,\"board\":{\"height\":21,\"width\":19,\"snakes\":[{\"id\":\"3e02eeed-67f2-4ded-8c7a-a554a5c8c588\",\"name\":\"local\",\"latency\":\"0\",\"health\":88,\"body\":[{\"x\":4,\"y\":19},{\"x\":4,\"y\":18},{\"x\":4,\"y\":17}],\"head\":{\"x\":4,\"y\":19},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}},{\"id\":\"4928d5c3-ed66-4e17-918e-6fa989a1c6c1\",\"name\":\"tiam\",\"latency\":\"0\",\"health\":88,\"body\":[{\"x\":14,\"y\":19},{\"x\":13,\"y\":19},{\"x\":12,\"y\":19}],\"head\":{\"x\":14,\"y\":19},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}],\"food\":[{\"x\":9,\"y\":11},{\"x\":14,\"y\":17},{\"x\":15,\"y\":11},{\"x\":9,\"y\":1},{\"x\":9,\"y\":17}],\"hazards\":[{\"x\":0,\"y\":20},{\"x\":2,\"y\":20},{\"x\":3,\"y\":20},{\"x\":4,\"y\":20},{\"x\":5,\"y\":20},{\"x\":6,\"y\":20},{\"x\":7,\"y\":20},{\"x\":8,\"y\":20},{\"x\":9,\"y\":20},{\"x\":10,\"y\":20},{\"x\":11,\"y\":20},{\"x\":12,\"y\":20},{\"x\":13,\"y\":20},{\"x\":14,\"y\":20},{\"x\":15,\"y\":20},{\"x\":16,\"y\":20},{\"x\":18,\"y\":20},{\"x\":0,\"y\":19},{\"x\":9,\"y\":19},{\"x\":18,\"y\":19},{\"x\":0,\"y\":18},{\"x\":2,\"y\":18},{\"x\":3,\"y\":18},{\"x\":5,\"y\":18},{\"x\":6,\"y\":18},{\"x\":7,\"y\":18},{\"x\":9,\"y\":18},{\"x\":11,\"y\":18},{\"x\":12,\"y\":18},{\"x\":13,\"y\":18},{\"x\":15,\"y\":18},{\"x\":16,\"y\":18},{\"x\":18,\"y\":18},{\"x\":0,\"y\":17},{\"x\":18,\"y\":17},{\"x\":0,\"y\":16},{\"x\":2,\"y\":16},{\"x\":3,\"y\":16},{\"x\":5,\"y\":16},{\"x\":7,\"y\":16},{\"x\":8,\"y\":16},{\"x\":9,\"y\":16},{\"x\":10,\"y\":16},{\"x\":11,\"y\":16},{\"x\":13,\"y\":16},{\"x\":15,\"y\":16},{\"x\":16,\"y\":16},{\"x\":18,\"y\":16},{\"x\":0,\"y\":15},{\"x\":5,\"y\":15},{\"x\":9,\"y\":15},{\"x\":13,\"y\":15},{\"x\":18,\"y\":15},{\"x\":0,\"y\":14},{\"x\":3,\"y\":14},{\"x\":5,\"y\":14},{\"x\":6,\"y\":14},{\"x\":7,\"y\":14},{\"x\":9,\"y\":14},{\"x\":11,\"y\":14},{\"x\":12,\"y\":14},{\"x\":13,\"y\":14},{\"x\":15,\"y\":14},{\"x\":18,\"y\":14},{\"x\":0,\"y\":13},{\"x\":3,\"y\":13},{\"x\":5,\"y\":13},{\"x\":13,\"y\":13},{\"x\":15,\"y\":13},{\"x\":18,\"y\":13},{\"x\":0,\"y\":12},{\"x\":1,\"y\":12},{\"x\":2,\"y\":12},{\"x\":3,\"y\":12},{\"x\":5,\"y\":12},{\"x\":7,\"y\":12},{\"x\":9,\"y\":12},{\"x\":11,\"y\":12},{\"x\":13,\"y\":12},{\"x\":15,\"y\":12},{\"x\":16,\"y\":12},{\"x\":17,\"y\":12},{\"x\":18,\"y\":12},{\"x\":7,\"y\":11},{\"x\":11,\"y\":11},{\"x\":0,\"y\":10},{\"x\":1,\"y\":10},{\"x\":2,\"y\":10},{\"x\":3,\"y\":10},{\"x\":5,\"y\":10},{\"x\":7,\"y\":10},{\"x\":9,\"y\":10},{\"x\":11,\"y\":10},{\"x\":13,\"y\":10},{\"x\":15,\"y\":10},{\"x\":16,\"y\":10},{\"x\":17,\"y\":10},{\"x\":18,\"y\":10},{\"x\":0,\"y\":9},{\"x\":3,\"y\":9},{\"x\":5,\"y\":9},{\"x\":13,\"y\":9},{\"x\":15,\"y\":9},{\"x\":18,\"y\":9},{\"x\":0,\"y\":8},{\"x\":3,\"y\":8},{\"x\":5,\"y\":8},{\"x\":7,\"y\":8},{\"x\":8,\"y\":8},{\"x\":9,\"y\":8},{\"x\":10,\"y\":8},{\"x\":11,\"y\":8},{\"x\":13,\"y\":8},{\"x\":15,\"y\":8},{\"x\":18,\"y\":8},{\"x\":0,\"y\":7},{\"x\":9,\"y\":7},{\"x\":18,\"y\":7},{\"x\":0,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":7,\"y\":6},{\"x\":9,\"y\":6},{\"x\":11,\"y\":6},{\"x\":12,\"y\":6},{\"x\":13,\"y\":6},{\"x\":15,\"y\":6},{\"x\":16,\"y\":6},{\"x\":18,\"y\":6},{\"x\":0,\"y\":5},{\"x\":3,\"y\":5},{\"x\":15,\"y\":5},{\"x\":18,\"y\":5},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":3,\"y\":4},{\"x\":5,\"y\":4},{\"x\":7,\"y\":4},{\"x\":8,\"y\":4},{\"x\":9,\"y\":4},{\"x\":10,\"y\":4},{\"x\":11,\"y\":4},{\"x\":13,\"y\":4},{\"x\":15,\"y\":4},{\"x\":17,\"y\":4},{\"x\":18,\"y\":4},{\"x\":0,\"y\":3},{\"x\":5,\"y\":3},{\"x\":9,\"y\":3},{\"x\":13,\"y\":3},{\"x\":18,\"y\":3},{\"x\":0,\"y\":2},{\"x\":2,\"y\":2},{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":6,\"y\":2},{\"x\":7,\"y\":2},{\"x\":9,\"y\":2},{\"x\":11,\"y\":2},{\"x\":12,\"y\":2},{\"x\":13,\"y\":2},{\"x\":14,\"y\":2},{\"x\":15,\"y\":2},{\"x\":16,\"y\":2},{\"x\":18,\"y\":2},{\"x\":0,\"y\":1},{\"x\":18,\"y\":1},{\"x\":0,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":6,\"y\":0},{\"x\":7,\"y\":0},{\"x\":8,\"y\":0},{\"x\":9,\"y\":0},{\"x\":10,\"y\":0},{\"x\":11,\"y\":0},{\"x\":12,\"y\":0},{\"x\":13,\"y\":0},{\"x\":14,\"y\":0},{\"x\":15,\"y\":0},{\"x\":16,\"y\":0},{\"x\":18,\"y\":0}]},\"you\":{\"id\":\"3e02eeed-67f2-4ded-8c7a-a554a5c8c588\",\"name\":\"local\",\"latency\":\"0\",\"health\":88,\"body\":[{\"x\":4,\"y\":19},{\"x\":4,\"y\":18},{\"x\":4,\"y\":17}],\"head\":{\"x\":4,\"y\":19},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}}")

	state := GameState{}
	_ = json.Unmarshal(g, &state)

	b := BuildBoard(state)

	if _, ok := b.Heads[MeId]; !ok {
		panic("Did not add me to board!")
	}
}
