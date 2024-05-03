package game

import (
	"fmt"
	"testing"
)

func TestTBMovement(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// _ _ h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := GameState{
		Turn: 0,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	if board.List[2].IsTripleStack() != true {
		panic("did not start with triple stack")
	}

	board.AdvanceBoardTB(SnakeMove{id, Up})

	if board.List[5].IsSnakeHead() != true {
		board.Print()
		panic("Did not move head up!")
	}

	if board.List[5].GetIdx() != 2 {
		board.Print()
		panic("Did not set tail idx on head")
	}

	if board.List[2].IsDoubleStack() != true {
		board.Print()
		panic("did not set double stack!")
	}

	if board.List[2].GetIdx() != 5 {
		board.Print()
		panic("Did not set head idx on tail")
	}

	board.AdvanceBoardTB(SnakeMove{id, Left})

	if board.List[4].IsSnakeHead() != true {
		board.Print()
		panic("Did not move head left!")
	}

	if board.List[4].GetIdx() != 2 {
		board.Print()
		panic("Did not set tail idx on head")
	}
	if board.List[5].GetIdx() != 4 {
		board.Print()
		panic("Did not set next idx on body")
	}
	if board.List[2].GetIdx() != 5 {
		board.Print()
		panic("Did not set next idx on tail")
	}

	if board.List[2].IsDoubleStack() || board.List[2].IsTripleStack() {
		panic("should not be stacked!")
	}

	board.AdvanceBoardTB(SnakeMove{id, Down})

	if board.List[1].IsSnakeHead() != true {
		panic("Did not move snake down!")
	}
	if board.List[2].IsSnakeSegment() != false {
		fmt.Println(board.List[2])
		panic("Did not move tail!")
	}

	if board.List[1].GetIdx() != 5 {
		board.Print()
		panic("Did not set tail idx on head")
	}
	if board.List[5].GetIdx() != 4 {
		board.Print()
		panic("Did not set next idx on body")
	}
	if board.List[4].GetIdx() != 1 {
		board.Print()
		panic("Did not set next idx on tail")
	}

	board.AdvanceBoardTB(SnakeMove{id, Right})

	if board.List[2].IsSnakeHead() != true {
		panic("Did not move snake right!")
	}
	if board.List[5].IsSnakeSegment() != false {
		panic("Did not move tail!")
	}

	if board.List[2].GetIdx() != 4 {
		board.Print()
		panic("Did not set tail idx on head")
	}
	if board.List[4].GetIdx() != 1 {
		board.Print()
		panic("Did not set next idx on body")
	}
	if board.List[1].GetIdx() != 2 {
		board.Print()
		panic("Did not set next idx on tail")
	}
}

func TestTBMovementTailFollow(t *testing.T) {
	//t.Skip()
	// _ _ _
	// s h _
	// s s _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Left})

	if !board.List[3].IsSnakeHead() {
		board.Print()
		panic("Did not move head")
	}

	if board.List[3].GetIdx() != 0 {
		board.Print()
		panic("head does not point to tail")
	}

	if board.List[4].GetIdx() != 3 {
		panic("body does not point to head")
	}

	if board.List[1].GetIdx() != 4 {
		panic("body does not point to next")
	}

	if board.List[0].GetIdx() != 1 {
		panic("tail does not point to next")
	}
}

func TestTBMovementThroughTails(t *testing.T) {
	//t.Skip()
	// o o e
	// o h _
	// s s _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Left})

	if board.List[3].id != id || board.List[3].tailId != enemyId {
		panic("Did not properly set head-tail")
	}

	if board.getTBSnakeIdx(3, id) != 1 {
		panic("Did not properly set head-tail next")
	}
	if board.getTBSnakeIdx(3, enemyId) != 6 {
		panic("Did not properly set head-tail next")
	}

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	if !board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.List[3].IsSnakeHead() || board.List[3].tailId != 0 {
		board.Print()
		panic("Did not properly remove head-tail")
	}

	if board.getTBSnakeIdx(3, id) != 1 {
		panic("Did not properly set head-tail next")
	}

	if board.List[1].GetIdx() != 4 {
		panic("Did not update tail to point to neck")
	}
	if board.List[4].GetIdx() != 3 {
		panic("Did not update neck to point to head")
	}

	if board.List[3].tailIdx != 0 || board.List[3].tailId != 0 {
		panic("Did not properly set enemy head-tail next")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

	if board.List[5].GetIdx() != 6 {
		panic("did not update head with new tail!")
	}
	if board.List[6].GetIdx() != 7 {
		panic("did not update tail with new next!")
	}

	if board.List[5].IsHeadTail() {
		board.Print()
		panic("Head should not also be tail!")
	}

}

func TestTBSideCollision(t *testing.T) {
	//t.Skip()
	// o o e
	// o h _
	// s s _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	//moves[enemyId] = Down
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

	if board.List[5].IsHeadTail() {
		board.Print()
		panic("Head should not also be tail!")
	}
}

func TestTBHeadCollision(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
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
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

	if board.List[5].IsHeadTail() {
		board.Print()
		panic("Head should not also be tail!")
	}
}

func TestEnemyTurnHeadCollisionDeath(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
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
		Turn: 3,
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

	//moves[enemyId] = Down
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

	if board.List[5].IsHeadTail() {
		board.Print()
		panic("Head should not also be tail!")
	}
}

func TestMyTurnHeadCollisionLife(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state := GameState{
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	if !board.IsSnakeAlive(id) {
		board.Print()
		panic("Should not remove me from board!")
	}

	if board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should have removed enemy!")
	}
}

func TestEnemyTurnHeadCollisionLife(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
	// s s h
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state := GameState{
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	if !board.IsSnakeAlive(id) {
		board.Print()
		panic("SHoudl not remove me from board!")
	}

	if board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should have removed enemy!")
	}
}

func TestTurnTailFollowFood(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ f
	// h s s
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []Coord{{2, 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	board2 := board.Clone()

	id := board.ids["me"]
	enemyId := board.ids["enemy"]

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Should remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

	board2.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})
	board2.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	if board2.IsSnakeAlive(id) {
		board2.Print()
		panic("Should remove me from board2!")
	}

	if !board2.IsSnakeAlive(enemyId) {
		board2.Print()
		panic("Should not have removed enemy!")
	}
}

func TestMyTurnOutofBounds(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
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
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Right})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}
}

func TestMyTurnOutofBoundsDown(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
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
		Turn: 3,
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

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Down})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}
}

func TestTBTailFollow(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ _
	// h s s
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
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

	fmt.Println("#####################################################################")
	board.Print()
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})
	board.Print()
	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Right})
	board.Print()

	if !board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

}

func TestTBEating(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ h f
	// _ s _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 0}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []Coord{{2, 1}},
		},
		You: me,
	}
	board := BuildBoard(state)

	id := board.ids["me"]

	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Right})

	if board.List[4].IsSnakeHead() {
		board.Print()
		panic("did set neck")
	}

	board.Print()
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})
	if board.List[1].IsDoubleStack() {
		board.Print()
		panic("did not remove double stack on move")
	}
}

func TestTBTailFollowFood(t *testing.T) {
	//t.Skip()
	// o o e
	// o _ f
	// h s s
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []Coord{{2, 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	board2 := board.Clone()

	id := board.ids["me"]
	enemyId := board.ids["enemy"]

	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}

	board2.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})
	board2.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Down})

	if board2.IsSnakeAlive(id) {
		board2.Print()
		panic("Did not remove me from board2!")
	}

	if !board2.IsSnakeAlive(enemyId) {
		board2.Print()
		panic("Should not have removed enemy!")
	}
}

func TestTBAdvanceBoardCrazyStuff(t *testing.T) {
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
		Turn: 4,
		You:  zero,
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

	board.AdvanceBoardTB(SnakeMove{id0, Up})
	board.AdvanceBoardTB(SnakeMove{id1, Right})
	board.AdvanceBoardTB(SnakeMove{id2, Left})
	board.AdvanceBoardTB(SnakeMove{id3, Up})

	if board.Healths[id1] > 0 {
		panic("Should not have killed id1!")
	}

	if board.Healths[id3] > 0 {
		panic("Should not have killed id3!")
	}

	snakeId, ok := board.List[9].GetSnakeId()
	if board.List[9].IsSnakeSegment() != true || !ok || snakeId != id0 {
		panic("lost snake id0 neck!")
	}
}

func TestTBAdvanceBoardOutOfBounds(t *testing.T) {
	//t.Skip()
	// right side
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
		Turn: 4,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]

	board.AdvanceBoardTB(SnakeMove{id, Right})

	if board.List[2].IsSnakeSegment() != false && board.List[1].IsSnakeSegment() != false && board.List[0].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// left side
	// _ _ _
	// _ _ _
	// h s s
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 0, Y: 0},
		Body:   []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	state = GameState{
		Turn: 4,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.ids["me"]

	board.AdvanceBoardTB(SnakeMove{id, Left})

	if board.List[2].IsSnakeSegment() != false && board.List[1].IsSnakeSegment() != false && board.List[0].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// top side
	// h _ _
	// s _ _
	// s _ _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 0, Y: 2},
		Body:   []Coord{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
	}
	state = GameState{
		Turn: 4,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.ids["me"]

	board.AdvanceBoardTB(SnakeMove{id, Up})

	if board.List[0].IsSnakeSegment() != false && board.List[3].IsSnakeSegment() != false && board.List[6].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// bottom side
	// _ _ s
	// _ _ s
	// _ _ h
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
	}
	state = GameState{
		Turn: 4,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.ids["me"]

	board.AdvanceBoardTB(SnakeMove{id, Down})

	if board.List[2].IsSnakeSegment() != false && board.List[5].IsSnakeSegment() != false && board.List[8].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}
}

func TestTBEatDoesNotRuinState(t *testing.T) {
	//t.Skip()
	// ff  __  ff  __  __  ff  __  ff  __  __  __
	// __  __  3s  3s  3h  3s  __  __  __  __  __
	// __  __  __  3s  ff  3s  __  __  __  ff  ff
	// __  __  __  3s  3s  3s  __  __  __  __  __
	// __  __  __  __  __  __  __  __  __  __  __
	// ff  1h  __  __  __  ff  __  ff  __  __  __
	// 1s  1s  1s  1s  __  __  __  __  __  __  __
	// 1s  1s  1s  1s  1s  1s  __  __  __  ff  ff
	// __  1s  1s  __  __  1s  ff  __  __  __  __
	// __  __  __  __  __  __  __  __  __  __  __
	// __  __  __  __  __  __  __  __  __  __  __
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []Coord{
			{1, 5},
			{1, 4},
			{0, 4},
			{0, 3},
			{1, 3},
			{1, 2},
			{2, 2},
			{2, 3},
			{2, 4},
			{3, 4},
			{3, 3},
			{4, 3},
			{5, 3},
			{5, 2},
		},
	}
	two := Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []Coord{
			{4, 9},
			{5, 9},
			{5, 8},
			{5, 7},
			{4, 7},
			{3, 7},
			{3, 8},
			{3, 9},
			{2, 9},
		},
	}
	state := GameState{
		Turn: 0,
		Board: Board{
			Snakes: []Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food: []Coord{
				{0, 5},
			},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	twoId := board.ids["two"]
	//threeId := board.ids["three"]

	board.AdvanceBoardTB(SnakeMove{id, Left})
	board.Print()
	if board.Healths[id] < 0 || board.Healths[twoId] < 0 {
		fmt.Println(board)
		panic("game did not end!")
	}
}
func TestTBRandomRollout(t *testing.T) {
	//t.Skip()
	// f _ f _ _ f _ f _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ e _ _ _ _ _ f f
	// _ f _ _ _ _ f _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// f _ f _ _ f _ f _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ f _ _ _ _ f f
	// e _ _ _ _ _ f _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ h _ _ _ _ _ _ _ _
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
	three := Battlesnake{
		ID:     "three",
		Health: 100,
		Body:   []Coord{{X: 3, Y: 8}, {X: 3, Y: 8}, {X: 3, Y: 8}},
	}
	state := GameState{
		Turn: 0,
		Board: Board{
			Snakes: []Battlesnake{me, two, three},
			Height: 11,
			Width:  11,
			Food: []Coord{
				{0, 10},
				{2, 10},
				{5, 10},
				{7, 10},
				{9, 8},
				{10, 8},
				{1, 7},
				{5, 7},
				{0, 5},
				{2, 5},
				{5, 5},
				{7, 5},
				{4, 3},
				{9, 3},
				{10, 3},
				{6, 2},
			},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	twoId := board.ids["two"]
	//threeId := board.ids["three"]

	fmt.Println("Running random rollout")
	board.RandomRolloutTB()

	fmt.Println(board)
	if board.Healths[id] > 1 && board.Healths[twoId] > 1 {
		fmt.Println(board)
		panic("game did not end!")
	}
}

func TestEnemyTurnHeadCollisionDeath2(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// _ _ _ _ _
	// _ _ _ e o
	// _ _ h s o
	// _ _ _ s o
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 3, Y: 2}, {X: 4, Y: 2}, {X: 4, Y: 1}, {X: 4, Y: 0}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 5,
			Width:  5,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.ids["me"]
	enemyId := board.ids["enemy"]

	board.Print()
	board.AdvanceBoardTB(SnakeMove{Id: id, Dir: Up})

	board.Print()
	board.AdvanceBoardTB(SnakeMove{Id: enemyId, Dir: Left})

	fmt.Println("222222222")
	board.Print()
	if board.IsSnakeAlive(id) {
		board.Print()
		panic("Did not remove me from board!")
	}

	if !board.IsSnakeAlive(enemyId) {
		board.Print()
		panic("Should not have removed enemy!")
	}
}

func TestTBGetNeighbors(t *testing.T) {
	//t.Skip()
	// o o e
	// o h _
	// s s _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)

	n := board.GetNeighborsTB(4)

	if len(n) != 2 {
		board.Print()
		fmt.Println(n)
		panic("thought you could go into snake")
	}

	// _ _ _ _ _
	// _ _ _ _ _
	// s s _ _ _
	// h s s s _
	// _ _ _ _ _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []Coord{
			{0, 1},
			{0, 2},
			{1, 2},
			{1, 1},
			{2, 1},
			{3, 1},
		},
	}
	state = GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 5,
			Width:  5,
		},
		You: me,
	}
	board = BuildBoard(state)

	n = board.GetNeighborsTB(5)

	if len(n) != 1 {
		board.Print()
		fmt.Println(n)
		panic("thought you could go into snake")
	}
}

func TestTBGetNeighborsStartingState(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ h _
	// _ _ _
	me := Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}},
	}
	state := GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)

	n := board.GetNeighborsTB(4)

	if len(n) != 4 {
		board.Print()
		fmt.Println(n)
		panic("Should be able to go in every direction")
	}

	//t.Skip()
	// _ _ _
	// _ h _
	// _ s _
	me = Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 0}},
	}
	state = GameState{
		Turn: 3,
		Board: Board{
			Snakes: []Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)

	n = board.GetNeighborsTB(4)

	if len(n) != 3 {
		board.Print()
		fmt.Println(n)
		panic("Should not be able to go down")
	}

}
