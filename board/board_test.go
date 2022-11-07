package board

import (
	"encoding/json"
	"fmt"
	b "github.com/nosnaws/tiam/battlesnake"
	"testing"
)

func createIslandBridgesGame() b.GameState {

	islandBridgesHz := []b.Coord{
		{X: 0, Y: 1},
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 4, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 1},
		{X: 6, Y: 0},
		{X: 9, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 1},
		{X: 10, Y: 4},
		{X: 9, Y: 5},
		{X: 10, Y: 5},
		{X: 10, Y: 6},
		{X: 10, Y: 9},
		{X: 10, Y: 10},
		{X: 9, Y: 10},
		{X: 6, Y: 10},
		{X: 5, Y: 10},
		{X: 5, Y: 9},
		{X: 4, Y: 10},
		{X: 1, Y: 10},
		{X: 0, Y: 10},
		{X: 0, Y: 9},
		{X: 0, Y: 6},
		{X: 0, Y: 5},
		{X: 1, Y: 5},
		{X: 0, Y: 4},
		{X: 2, Y: 4},
		{X: 3, Y: 5},
		{X: 4, Y: 5},
		{X: 5, Y: 5},
		{X: 6, Y: 5},
		{X: 7, Y: 5},
		{X: 5, Y: 7},
		{X: 5, Y: 6},
		{X: 5, Y: 5},
		{X: 5, Y: 4},
		{X: 5, Y: 3},
	}
	state := b.GameState{
		Turn: 0,
		Board: b.Board{
			Height:  11,
			Width:   11,
			Hazards: islandBridgesHz,
		},
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
				Settings: b.Settings{
					HazardDamagePerTurn: 100,
				},
			},
		},
	}

	return state
}

func TestCantMoveBackOnNeck(t *testing.T) {
	//t.Skip()
	// x x _ _ x x x _ e x x
	// x _ _ _ _ x _ s s _ x
	// _ _ _ _ _ _ _ _ _ _ _
	// _ _ _ _ _ x _ _ _ _ _
	// x _ _ _ _ x _ _ _ _ x
	// x x _ x x x x x _ x x
	// x _ _ _ _ x _ _ _ _ x
	// _ _ _ _ _ x _ _ _ _ _
	// _ _ _ _ _ _ _ _ _ _ _
	// x _ _ _ _ x _ s s _ x
	// x x _ _ x x x _ h x x
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []b.Coord{
			{X: 8, Y: 0},
			{X: 8, Y: 1},
			{X: 7, Y: 1},
		},
	}
	two := b.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []b.Coord{
			{X: 8, Y: 10},
			{X: 8, Y: 9},
			{X: 7, Y: 9},
		},
	}
	state := createIslandBridgesGame()
	state.Board.Snakes = []b.Battlesnake{me, two}
	board := BuildBoard(state)
	cart := GetCartesianProductOfMoves(&board)
	//moves := board.GetMovesForSnake(MeId)

	if len(cart) != 2 {
		board.Print()
		fmt.Println(cart)
		panic("Should only have 1 moves!")
	}

}

func TestGetMovesForSnake(t *testing.T) {
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
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body: []b.Coord{
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
	two := b.Battlesnake{
		ID:     "two",
		Health: 100,
		Body: []b.Coord{
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
	state := b.GameState{
		Turn: 0,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, two},
			Height: 11,
			Width:  11,
			Food: []b.Coord{
				{X: 7, Y: 9},
				{X: 10, Y: 9},
			},
		},
		You: me,
	}
	board := BuildBoard(state)
	moves := board.GetMovesForSnake(1)

	if len(moves) != 3 {
		board.Print()
		fmt.Println(moves)
		panic("Should only have 3 moves!")
	}

}

func TestHeadToHead(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// _ s s s h
	// s s s e f
	// _ _ _ _ _
	// _ _ _ _ _
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 4, Y: 3}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}},
	}
	two := b.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []b.Coord{{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state := b.GameState{
		Board: b.Board{
			Snakes: []b.Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food:   []b.Coord{{X: 4, Y: 2}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]
	twoId := board.Ids["two"]

	moves := make(map[SnakeId]Move)
	moves[id] = Down
	moves[twoId] = Right
	board.Print()
	board.AdvanceBoard(moves)

	fmt.Println("HERE")
	board.Print()
	if !board.IsGameOver() {
		board.Print()
		panic("game did not end!")
	}
}

func TestRandomSnakeCollsionWrapped(t *testing.T) {
	//t.Skip()
	// s s _
	// e f h
	// _ s s
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 1},
		Body:   []b.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}},
	}
	two := b.Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   b.Coord{X: 0, Y: 1},
		Body:   []b.Coord{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 2}},
	}
	state := b.GameState{
		Board: b.Board{
			Snakes: []b.Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 1, Y: 1}},
		},
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]
	twoId := board.Ids["two"]

	moves := make(map[SnakeId]Move)
	moves[id] = Right
	moves[twoId] = Down
	board.AdvanceBoard(moves)

	if board.Healths[id] > 1 {
		board.Print()
		panic("game did not end!")
	}
}

func TestGetNeighbors(t *testing.T) {
	//t.Skip()
	// e _ _
	// s s _
	// h s s
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	enemy := b.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []b.Coord{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 1, Y: 1}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, enemy},
			Height: 11,
			Width:  11,
		},
		You: me,
	}
	board := BuildBoard(state)
	//id := board.ids["me"]

	moves := board.GetNeighbors(0)

	if len(moves) != 0 {
		fmt.Println(moves)
		board.Print()
		panic("Should not be able to move!")
	}
}

func TestGetSnakeMoves(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := board.GetMovesForSnake(id)

	if len(moves) != 1 {
		fmt.Println(moves)
		board.Print()
		panic("Should only be able to move up!")
	}

	// wrapped!
	// _ _ _
	// _ _ _
	// s s h
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state = b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = board.GetMovesForSnake(id)

	if len(moves) != 3 {
		fmt.Println(moves)
		panic("Should be able to go up,down, and right!")
	}

	// wrapped with snake eating
	// f e s
	// s _ s
	// s h _
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 1, Y: 0},
		Body:   []b.Coord{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	}
	two := b.Battlesnake{
		ID:     "two",
		Health: 100,
		Head:   b.Coord{X: 1, Y: 2},
		Body:   []b.Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
	}
	state = b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, two},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 0, Y: 2}},
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.Ids["me"]
	enemy := board.Ids["two"]

	moves2 := make(map[SnakeId]Move)
	moves2[id] = Right
	moves2[enemy] = Left

	board.AdvanceBoard(moves2)
	board.Print()
	moves = board.GetMovesForSnake(id)
	if !board.List[8].IsDoubleStack() {
		panic("should be double stack!")
	}

	if len(moves) != 2 {
		board.Print()
		fmt.Println(moves)
		panic("Should be able to go right and up!")
	}
}

func TestBoardCreationTurn0(t *testing.T) {
	//t.Skip()
	// Arrange
	me := b.Battlesnake{
		// Length 3, facing right
		Head: b.Coord{X: 2, Y: 0},
		Body: []b.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := b.GameState{
		Turn: 0,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 0}},
			Food:    []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}

	board := BuildBoard(state)

	if board.List[2].id != 1 {
		panic("YouId is not 1")
	}
	if board.List[2].IsTripleStack() != true {
		panic("Did not triple stack snake!")
	}

}

func TestBoardCreationTurn1(t *testing.T) {
	//t.Skip()
	// Arrange
	me := b.Battlesnake{
		// Length 3, facing right
		Head: b.Coord{X: 1, Y: 0},
		Body: []b.Coord{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := b.GameState{
		Turn: 1,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 0}},
			Food:    []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}

	board := BuildBoard(state)

	if board.List[2].id != 1 {
		panic("YouId is not 1")
	}
	if board.List[1].IsSnakeHead() != true {
		panic("did not place snake head!")
	}
	if board.List[2].IsDoubleStack() != true {
		panic("Did not double stack snake!")
	}

}

func TestBoardCreation(t *testing.T) {
	//t.Skip()
	// Arrange
	me := b.Battlesnake{
		// Length 3, facing right
		Head: b.Coord{X: 2, Y: 0},
		Body: []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 0}},
			Food:    []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}

	board := BuildBoard(state)

	if board.List[2].id != 1 {
		panic("YouId is not 1")
	}

	if board.List[2].IsHazard() != true {
		panic("Did not create hazard")
	}

	if board.List[5].IsFood() != true {
		panic("Did not create food")
	}

	head := board.List[2]
	tail := board.List[head.GetIdx()]
	current := tail
	for current != head {
		current = board.List[current.GetIdx()]
	}

	if current != head {
		panic("Snake does not loop to head!")
	}
}

func TestKill(t *testing.T) {
	//t.Skip()
	// Arrange
	me := b.Battlesnake{
		// Length 3, facing right
		ID:   "me",
		Head: b.Coord{X: 2, Y: 0},
		Body: []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 0}},
			Food:    []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}

	board := BuildBoard(state)
	id := board.Ids["me"]
	board.kill(id)

	if board.List[2].id == 1 {
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
	//t.Skip()
	// _ _ _
	// _ _ _
	// _ _ h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := b.GameState{
		Turn: 0,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]
	board.Print()

	if board.List[2].IsTripleStack() != true {
		panic("did not start with triple stack")
	}

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)
	board.Print()

	if board.List[5].IsSnakeHead() != true {
		fmt.Println(board)
		fmt.Println(board.List[5].idx)
		panic("Did not move snake up!")
	}
	if board.List[2].IsDoubleStack() != true {
		fmt.Println(board)
		panic("did not set double stack!")
	}

	moves = make(map[SnakeId]Move)
	moves[id] = Left
	board.AdvanceBoard(moves)
	board.Print()

	if board.List[4].IsSnakeHead() != true {
		fmt.Println(board)
		panic("Did not move snake left!")
	}
	if board.List[2].IsDoubleStack() || board.List[2].IsTripleStack() {
		panic("should not be stacked!")
	}

	moves = make(map[SnakeId]Move)
	moves[id] = Down
	board.AdvanceBoard(moves)
	board.Print()

	if board.List[1].IsSnakeHead() != true {
		panic("Did not move snake down!")
	}
	if board.List[2].IsSnakeSegment() != false {
		fmt.Println(board.List[2])
		panic("Did not move tail!")
	}

	moves = make(map[SnakeId]Move)
	moves[id] = Right
	board.AdvanceBoard(moves)
	board.Print()

	if board.List[2].IsSnakeHead() != true {
		panic("Did not move snake right!")
	}
	if board.List[5].IsSnakeSegment() != false {
		panic("Did not move tail!")
	}
}

func TestAdvanceBoardTurnDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// s s h
	me := b.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}

	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		fmt.Println(board.Healths[id])
		panic("Did not decrement health properly!")
	}
}

func TestAdvanceBoardHazardDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ x
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 50,
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Settings: b.Settings{
					HazardDamagePerTurn: 100,
				},
			},
		},
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] > 0 {
		fmt.Println(board.Healths[id])
		board.Print()
		panic("Snake did not die!")
	}
}

func TestAdvanceBoardStackedHazardDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ x
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 1}, {X: 2, Y: 1}, {X: 2, Y: 1}},
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Settings: b.Settings{
					HazardDamagePerTurn: 10,
				},
			},
		},
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] != 69 {
		panic("Snake did not take stacked hazard damage!")
	}
}

func TestAdvanceBoardHazardHealing(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ x
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Settings: b.Settings{
					HazardDamagePerTurn: -50,
				},
			},
		},
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		panic("Snake did not heal!")
	}
}

func TestAdvanceBoardHazardStacked(t *testing.T) {
	t.Skip()
	// _ _ _
	// _ _ x
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes:  []b.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []b.Coord{{X: 2, Y: 1}, {X: 2, Y: 1}},
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Settings: b.Settings{
					HazardDamagePerTurn: 25,
				},
			},
		},
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] != 0 {
		panic("Snake did not die!")
	}
}

func TestEatWhenStacked(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ f
	// _ s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.Print()
	board.AdvanceBoard(moves)
	board.Print()

	if board.Healths[id] != 100 {
		panic("Snake did not eat!")
	}
	if board.Lengths[id] != 4 {
		panic("Snake did not grow!")
	}
	if !board.List[1].IsSnakeSegment() {
		panic("Moved tail!")
	}
	if !board.List[1].IsDoubleStack() {
		panic("Did not set double stack!!")
	}

	moves = make(map[SnakeId]Move)
	moves[id] = Up
	board.Print()
	board.AdvanceBoard(moves)
	board.Print()

	if board.Healths[id] != 99 {
		panic("Did not reduce health!")
	}
	if !board.List[1].IsSnakeSegment() {
		panic("Moved tail!")
	}
}

func TestAdvanceBoardEatFood(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ f
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 50,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.Print()
	board.AdvanceBoard(moves)
	board.Print()

	if board.Healths[id] != 100 {
		panic("Snake did not eat!")
	}
	if board.Lengths[id] != 4 {
		panic("Snake did not grow!")
	}
	if board.List[0].IsSnakeSegment() {
		panic("Did not move tail")
	}

	moves = make(map[SnakeId]Move)
	moves[id] = Up
	board.Print()
	board.AdvanceBoard(moves)
	board.Print()

	if board.Healths[id] != 99 {
		panic("Did not reduce health!")
	}
	if !board.List[1].IsSnakeSegment() {
		panic("Moved tail!")
	}
}

func TestAdvanceBoardOutOfBounds(t *testing.T) {
	//t.Skip()
	// right side
	// _ _ _
	// _ _ _
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Right
	board.AdvanceBoard(moves)

	if board.List[2].IsSnakeSegment() != false && board.List[1].IsSnakeSegment() != false && board.List[0].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// left side
	// _ _ _
	// _ _ _
	// h s s
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 0, Y: 0},
		Body:   []b.Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	state = b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = make(map[SnakeId]Move)
	moves[id] = Left
	board.AdvanceBoard(moves)

	if board.List[2].IsSnakeSegment() != false && board.List[1].IsSnakeSegment() != false && board.List[0].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// top side
	// h _ _
	// s _ _
	// s _ _
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 0, Y: 2},
		Body:   []b.Coord{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
	}
	state = b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.List[0].IsSnakeSegment() != false && board.List[3].IsSnakeSegment() != false && board.List[6].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// bottom side
	// _ _ s
	// _ _ s
	// _ _ h
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   b.Coord{X: 2, Y: 0},
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
	}
	state = b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = make(map[SnakeId]Move)
	moves[id] = Down
	board.AdvanceBoard(moves)

	if board.List[2].IsSnakeSegment() != false && board.List[5].IsSnakeSegment() != false && board.List[8].IsSnakeSegment() != false {
		panic("Did not remove snake from board!")
	}
}

func TestAdvanceBoardHeadCollision(t *testing.T) {
	t.Skip()
	// o o e
	// o _ f
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := b.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]
	enemyId := board.Ids["enemy"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	moves[enemyId] = Down
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
	//t.Skip()
	// o o _
	// o e _
	// s s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := b.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []b.Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]
	enemyId := board.Ids["enemy"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	moves[enemyId] = Down
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
	//t.Skip()
	// _ s _
	// _ s h
	// _ s s
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Left
	board.AdvanceBoard(moves)

	if board.Healths[id] != 0 {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardFollowTail(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ s s
	// _ s h
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] < 1 {
		panic("Snake was killed!")
	}

	if board.List[5].IsSnakeHead() != true {
		panic("Did not move snake!")
	}

	// Follow other snake tail
	// _ e e
	// _ s e
	// _ s h
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	}
	enemy := b.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []b.Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
	}
	state = b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.Ids["me"]
	enemyId := board.Ids["enemy"]

	moves = make(map[SnakeId]Move)
	moves[id] = Up
	moves[enemyId] = Left
	fmt.Println("advancing board")
	board.AdvanceBoard(moves)
	fmt.Println("board advanced")

	if board.Healths[id] < 1 {
		panic("Snake was killed!")
	}

	if board.List[5].IsSnakeHead() != true {
		fmt.Println(board)
		panic("Did not move snake!")
	}

	// Follow other snake tail
	// f e e
	// _ s e
	// _ s h
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	}
	enemy = b.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []b.Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 1}},
	}
	state = b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []b.Coord{{X: 0, Y: 2}},
		},
		You: me,
	}
	board = BuildBoard(state)
	id = board.Ids["me"]
	enemyId = board.Ids["enemy"]

	moves = make(map[SnakeId]Move)
	moves[id] = Up
	moves[enemyId] = Left
	fmt.Println("advancing board")
	board.AdvanceBoard(moves)
	fmt.Println("board advanced")

	if board.Healths[id] != 0 {
		panic("Snake was not killed!")
	}

	if board.List[5].IsSnakeHead() != false {
		fmt.Println(board)
		panic("Did not remove snake!")
	}
}

func TestAdvanceBoardMoveOnNeck(t *testing.T) {
	//t.Skip()
	// _ s _
	// _ s h
	// _ s s
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state := b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Down
	board.AdvanceBoard(moves)

	if board.Healths[id] != 0 {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardWrapped(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ h
	// _ s s
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := make(map[SnakeId]Move)
	moves[id] = Right
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.List[3].IsSnakeHead() != true {
		panic("Did not wrapped to right!")
	}

	// _ _ h
	// _ _ s
	// _ _ s
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}},
	}
	state = b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = make(map[SnakeId]Move)
	moves[id] = Up
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.List[2].IsSnakeHead() != true {
		panic("Did not wrapped to bottom!")
	}

	// _ s _
	// _ s _
	// _ h _
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state = b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = make(map[SnakeId]Move)
	moves[id] = Down
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.List[7].IsSnakeHead() != true {
		panic("Did not wrapped to top!")
	}

	// _ _ _
	// s s h
	// _ _ _
	me = b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 1}},
	}
	state = b.GameState{
		Turn: 4,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = BuildBoard(state)
	id = board.Ids["me"]

	moves = make(map[SnakeId]Move)
	moves[id] = Right
	board.AdvanceBoard(moves)

	if board.Healths[id] != 99 {
		panic("Should not have killed me!")
	}

	if board.List[3].IsSnakeHead() != true {
		panic("Did not wrapped to left!")
	}
}

func TestSnakeEating(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// _ _ _ _ _
	// _ _ _ s s
	// _ _ _ h s
	// _ _ _ _ _
	// snake just ate
	me := b.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []b.Coord{{X: 3, Y: 1}, {X: 4, Y: 1}, {X: 4, Y: 2}, {X: 3, Y: 2}, {X: 3, Y: 2}},
	}
	state := b.GameState{
		Turn: 3,
		Board: b.Board{
			Snakes: []b.Battlesnake{me},
			Height: 5,
			Width:  5,
		},
		You: me,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := BuildBoard(state)
	id := board.Ids["me"]

	moves := board.GetMovesForSnake(id)

	if len(moves) != 2 {
		dirIndex := IndexInDirection(Up, board.Heads[id], 5, 5, true)
		fmt.Println(board.List[dirIndex].IsDoubleStack())
		fmt.Println(moves)
		panic("Should not go up!")
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
	zero := b.Battlesnake{
		ID:     "0",
		Health: 100,
		Body:   []b.Coord{{X: 9, Y: 1}, {X: 9, Y: 0}, {X: 10, Y: 0}},
	}
	one := b.Battlesnake{
		ID:     "1",
		Health: 100,
		Body:   []b.Coord{{X: 8, Y: 0}, {X: 8, Y: 10}, {X: 8, Y: 9}},
	}
	two := b.Battlesnake{
		ID:     "2",
		Health: 100,
		Body:   []b.Coord{{X: 4, Y: 4}, {X: 5, Y: 4}, {X: 5, Y: 5}},
	}
	three := b.Battlesnake{
		ID:     "3",
		Health: 100,
		Body:   []b.Coord{{X: 9, Y: 10}, {X: 10, Y: 10}, {X: 0, Y: 10}},
	}

	state := b.GameState{
		Board: b.Board{
			Snakes: []b.Battlesnake{zero, one, two, three},
			Height: 11,
			Width:  11,
		},
		Turn: 4,
		You:  zero,
		Game: b.Game{
			Ruleset: b.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := BuildBoard(state)
	id0 := board.Ids["0"]
	id1 := board.Ids["1"]
	id2 := board.Ids["2"]
	id3 := board.Ids["3"]

	moves := make(map[SnakeId]Move)
	moves[id0] = Up
	moves[id1] = Right
	moves[id2] = Left
	moves[id3] = Up
	board.AdvanceBoard(moves)

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
func TestBuildBoardSnakeAte(t *testing.T) {
	g := []byte("{\"game\":{\"id\":\"a8732cf3-42b6-4012-bbb5-d1cc4e03f397\",\"ruleset\":{\"name\":\"standard\",\"version\":\"cli\",\"settings\":{\"foodSpawnChance\":15,\"minimumFood\":1,\"hazardDamagePerTurn\":14,\"hazardMap\":\"\",\"hazardMapAuthor\":\"\",\"royale\":{\"shrinkEveryNTurns\":25},\"squad\":{\"allowBodyCollisions\":false,\"sharedElimination\":false,\"sharedHealth\":false,\"sharedLength\":false}}},\"map\":\"standard\",\"timeout\":500,\"source\":\"\"},\"turn\":2,\"board\":{\"height\":11,\"width\":11,\"snakes\":[{\"id\":\"45f8bf5b-02ea-4ccc-99f9-3079f5bbb805\",\"name\":\"tiam\",\"latency\":\"0\",\"health\":100,\"body\":[{\"x\":8,\"y\":10},{\"x\":8,\"y\":9},{\"x\":9,\"y\":9},{\"x\":9,\"y\":9}],\"head\":{\"x\":8,\"y\":10},\"length\":4,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}},{\"id\":\"2bd13899-e75e-4d8f-b6fb-fad018d04b7f\",\"name\":\"local\",\"latency\":\"0\",\"health\":98,\"body\":[{\"x\":10,\"y\":0},{\"x\":10,\"y\":1},{\"x\":9,\"y\":1}],\"head\":{\"x\":10,\"y\":0},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}],\"food\":[{\"x\":8,\"y\":0},{\"x\":5,\"y\":5},{\"x\":1,\"y\":2}],\"hazards\":[]},\"you\":{\"id\":\"2bd13899-e75e-4d8f-b6fb-fad018d04b7f\",\"name\":\"local\",\"latency\":\"0\",\"health\":98,\"body\":[{\"x\":10,\"y\":0},{\"x\":10,\"y\":1},{\"x\":9,\"y\":1}],\"head\":{\"x\":10,\"y\":0},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}}")

	state := b.GameState{}
	_ = json.Unmarshal(g, &state)

	b := BuildBoard(state)

	// x: 9, y: 9 = 9 * 11 + 9 = 108
	ind := pointToIndex(Point{X: 9, Y: 9}, b.Width)
	fmt.Println("is double", b.List[ind].IsDoubleStack())
	if !b.List[108].IsDoubleStack() {
		panic("Did not add double stack after eating!")
	}
}

func TestBuildBigBoard(t *testing.T) {
	g := []byte("{\"game\":{\"id\":\"c20df634-4097-471b-a558-c6e96ac56620\",\"ruleset\":{\"name\":\"wrapped\",\"version\":\"cli\",\"settings\":{\"foodSpawnChance\":15,\"minimumFood\":1,\"hazardDamagePerTurn\":100,\"hazardMap\":\"\",\"hazardMapAuthor\":\"\",\"royale\":{\"shrinkEveryNTurns\":25},\"squad\":{\"allowBodyCollisions\":false,\"sharedElimination\":false,\"sharedHealth\":false,\"sharedLength\":false}}},\"map\":\"arcade_maze\",\"timeout\":500,\"source\":\"\"},\"turn\":7,\"board\":{\"height\":21,\"width\":19,\"snakes\":[{\"id\":\"9eecc83c-283a-446c-98e0-b1f5e36d9c83\",\"name\":\"tiam\",\"latency\":\"0\",\"health\":93,\"body\":[{\"x\":10,\"y\":10},{\"x\":10,\"y\":9},{\"x\":11,\"y\":9}],\"head\":{\"x\":10,\"y\":10},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}},{\"id\":\"9cca4f6e-8872-42f1-abaf-685887e0697d\",\"name\":\"local\",\"latency\":\"0\",\"health\":93,\"body\":[{\"x\":4,\"y\":14},{\"x\":4,\"y\":13},{\"x\":4,\"y\":12}],\"head\":{\"x\":4,\"y\":14},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}],\"food\":[{\"x\":9,\"y\":11},{\"x\":4,\"y\":17}],\"hazards\":[{\"x\":0,\"y\":20},{\"x\":2,\"y\":20},{\"x\":3,\"y\":20},{\"x\":4,\"y\":20},{\"x\":5,\"y\":20},{\"x\":6,\"y\":20},{\"x\":7,\"y\":20},{\"x\":8,\"y\":20},{\"x\":9,\"y\":20},{\"x\":10,\"y\":20},{\"x\":11,\"y\":20},{\"x\":12,\"y\":20},{\"x\":13,\"y\":20},{\"x\":14,\"y\":20},{\"x\":15,\"y\":20},{\"x\":16,\"y\":20},{\"x\":18,\"y\":20},{\"x\":0,\"y\":19},{\"x\":9,\"y\":19},{\"x\":18,\"y\":19},{\"x\":0,\"y\":18},{\"x\":2,\"y\":18},{\"x\":3,\"y\":18},{\"x\":5,\"y\":18},{\"x\":6,\"y\":18},{\"x\":7,\"y\":18},{\"x\":9,\"y\":18},{\"x\":11,\"y\":18},{\"x\":12,\"y\":18},{\"x\":13,\"y\":18},{\"x\":15,\"y\":18},{\"x\":16,\"y\":18},{\"x\":18,\"y\":18},{\"x\":0,\"y\":17},{\"x\":18,\"y\":17},{\"x\":0,\"y\":16},{\"x\":2,\"y\":16},{\"x\":3,\"y\":16},{\"x\":5,\"y\":16},{\"x\":7,\"y\":16},{\"x\":8,\"y\":16},{\"x\":9,\"y\":16},{\"x\":10,\"y\":16},{\"x\":11,\"y\":16},{\"x\":13,\"y\":16},{\"x\":15,\"y\":16},{\"x\":16,\"y\":16},{\"x\":18,\"y\":16},{\"x\":0,\"y\":15},{\"x\":5,\"y\":15},{\"x\":9,\"y\":15},{\"x\":13,\"y\":15},{\"x\":18,\"y\":15},{\"x\":0,\"y\":14},{\"x\":3,\"y\":14},{\"x\":5,\"y\":14},{\"x\":6,\"y\":14},{\"x\":7,\"y\":14},{\"x\":9,\"y\":14},{\"x\":11,\"y\":14},{\"x\":12,\"y\":14},{\"x\":13,\"y\":14},{\"x\":15,\"y\":14},{\"x\":18,\"y\":14},{\"x\":0,\"y\":13},{\"x\":3,\"y\":13},{\"x\":5,\"y\":13},{\"x\":13,\"y\":13},{\"x\":15,\"y\":13},{\"x\":18,\"y\":13},{\"x\":0,\"y\":12},{\"x\":1,\"y\":12},{\"x\":2,\"y\":12},{\"x\":3,\"y\":12},{\"x\":5,\"y\":12},{\"x\":7,\"y\":12},{\"x\":9,\"y\":12},{\"x\":11,\"y\":12},{\"x\":13,\"y\":12},{\"x\":15,\"y\":12},{\"x\":16,\"y\":12},{\"x\":17,\"y\":12},{\"x\":18,\"y\":12},{\"x\":7,\"y\":11},{\"x\":11,\"y\":11},{\"x\":0,\"y\":10},{\"x\":1,\"y\":10},{\"x\":2,\"y\":10},{\"x\":3,\"y\":10},{\"x\":5,\"y\":10},{\"x\":7,\"y\":10},{\"x\":9,\"y\":10},{\"x\":11,\"y\":10},{\"x\":13,\"y\":10},{\"x\":15,\"y\":10},{\"x\":16,\"y\":10},{\"x\":17,\"y\":10},{\"x\":18,\"y\":10},{\"x\":0,\"y\":9},{\"x\":3,\"y\":9},{\"x\":5,\"y\":9},{\"x\":13,\"y\":9},{\"x\":15,\"y\":9},{\"x\":18,\"y\":9},{\"x\":0,\"y\":8},{\"x\":3,\"y\":8},{\"x\":5,\"y\":8},{\"x\":7,\"y\":8},{\"x\":8,\"y\":8},{\"x\":9,\"y\":8},{\"x\":10,\"y\":8},{\"x\":11,\"y\":8},{\"x\":13,\"y\":8},{\"x\":15,\"y\":8},{\"x\":18,\"y\":8},{\"x\":0,\"y\":7},{\"x\":9,\"y\":7},{\"x\":18,\"y\":7},{\"x\":0,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":7,\"y\":6},{\"x\":9,\"y\":6},{\"x\":11,\"y\":6},{\"x\":12,\"y\":6},{\"x\":13,\"y\":6},{\"x\":15,\"y\":6},{\"x\":16,\"y\":6},{\"x\":18,\"y\":6},{\"x\":0,\"y\":5},{\"x\":3,\"y\":5},{\"x\":15,\"y\":5},{\"x\":18,\"y\":5},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":3,\"y\":4},{\"x\":5,\"y\":4},{\"x\":7,\"y\":4},{\"x\":8,\"y\":4},{\"x\":9,\"y\":4},{\"x\":10,\"y\":4},{\"x\":11,\"y\":4},{\"x\":13,\"y\":4},{\"x\":15,\"y\":4},{\"x\":17,\"y\":4},{\"x\":18,\"y\":4},{\"x\":0,\"y\":3},{\"x\":5,\"y\":3},{\"x\":9,\"y\":3},{\"x\":13,\"y\":3},{\"x\":18,\"y\":3},{\"x\":0,\"y\":2},{\"x\":2,\"y\":2},{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":6,\"y\":2},{\"x\":7,\"y\":2},{\"x\":9,\"y\":2},{\"x\":11,\"y\":2},{\"x\":12,\"y\":2},{\"x\":13,\"y\":2},{\"x\":14,\"y\":2},{\"x\":15,\"y\":2},{\"x\":16,\"y\":2},{\"x\":18,\"y\":2},{\"x\":0,\"y\":1},{\"x\":18,\"y\":1},{\"x\":0,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":6,\"y\":0},{\"x\":7,\"y\":0},{\"x\":8,\"y\":0},{\"x\":9,\"y\":0},{\"x\":10,\"y\":0},{\"x\":11,\"y\":0},{\"x\":12,\"y\":0},{\"x\":13,\"y\":0},{\"x\":14,\"y\":0},{\"x\":15,\"y\":0},{\"x\":16,\"y\":0},{\"x\":18,\"y\":0}]},\"you\":{\"id\":\"9cca4f6e-8872-42f1-abaf-685887e0697d\",\"name\":\"local\",\"latency\":\"0\",\"health\":93,\"body\":[{\"x\":4,\"y\":14},{\"x\":4,\"y\":13},{\"x\":4,\"y\":12}],\"head\":{\"x\":4,\"y\":14},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}}")

	state := b.GameState{}
	_ = json.Unmarshal(g, &state)

	b := BuildBoard(state)

	if _, ok := b.Heads[MeId]; !ok {
		panic("Did not add me to board!")
	}
}

func TestBuildBigBoard2(t *testing.T) {
	g := []byte("{\"game\":{\"id\":\"676d530f-2d7f-4590-8fda-d15579e5b6fe\",\"ruleset\":{\"name\":\"wrapped\",\"version\":\"cli\",\"settings\":{\"foodSpawnChance\":15,\"minimumFood\":1,\"hazardDamagePerTurn\":100,\"hazardMap\":\"\",\"hazardMapAuthor\":\"\",\"royale\":{\"shrinkEveryNTurns\":25},\"squad\":{\"allowBodyCollisions\":false,\"sharedElimination\":false,\"sharedHealth\":false,\"sharedLength\":false}}},\"map\":\"arcade_maze\",\"timeout\":500,\"source\":\"\"},\"turn\":12,\"board\":{\"height\":21,\"width\":19,\"snakes\":[{\"id\":\"3e02eeed-67f2-4ded-8c7a-a554a5c8c588\",\"name\":\"local\",\"latency\":\"0\",\"health\":88,\"body\":[{\"x\":4,\"y\":19},{\"x\":4,\"y\":18},{\"x\":4,\"y\":17}],\"head\":{\"x\":4,\"y\":19},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}},{\"id\":\"4928d5c3-ed66-4e17-918e-6fa989a1c6c1\",\"name\":\"tiam\",\"latency\":\"0\",\"health\":88,\"body\":[{\"x\":14,\"y\":19},{\"x\":13,\"y\":19},{\"x\":12,\"y\":19}],\"head\":{\"x\":14,\"y\":19},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}],\"food\":[{\"x\":9,\"y\":11},{\"x\":14,\"y\":17},{\"x\":15,\"y\":11},{\"x\":9,\"y\":1},{\"x\":9,\"y\":17}],\"hazards\":[{\"x\":0,\"y\":20},{\"x\":2,\"y\":20},{\"x\":3,\"y\":20},{\"x\":4,\"y\":20},{\"x\":5,\"y\":20},{\"x\":6,\"y\":20},{\"x\":7,\"y\":20},{\"x\":8,\"y\":20},{\"x\":9,\"y\":20},{\"x\":10,\"y\":20},{\"x\":11,\"y\":20},{\"x\":12,\"y\":20},{\"x\":13,\"y\":20},{\"x\":14,\"y\":20},{\"x\":15,\"y\":20},{\"x\":16,\"y\":20},{\"x\":18,\"y\":20},{\"x\":0,\"y\":19},{\"x\":9,\"y\":19},{\"x\":18,\"y\":19},{\"x\":0,\"y\":18},{\"x\":2,\"y\":18},{\"x\":3,\"y\":18},{\"x\":5,\"y\":18},{\"x\":6,\"y\":18},{\"x\":7,\"y\":18},{\"x\":9,\"y\":18},{\"x\":11,\"y\":18},{\"x\":12,\"y\":18},{\"x\":13,\"y\":18},{\"x\":15,\"y\":18},{\"x\":16,\"y\":18},{\"x\":18,\"y\":18},{\"x\":0,\"y\":17},{\"x\":18,\"y\":17},{\"x\":0,\"y\":16},{\"x\":2,\"y\":16},{\"x\":3,\"y\":16},{\"x\":5,\"y\":16},{\"x\":7,\"y\":16},{\"x\":8,\"y\":16},{\"x\":9,\"y\":16},{\"x\":10,\"y\":16},{\"x\":11,\"y\":16},{\"x\":13,\"y\":16},{\"x\":15,\"y\":16},{\"x\":16,\"y\":16},{\"x\":18,\"y\":16},{\"x\":0,\"y\":15},{\"x\":5,\"y\":15},{\"x\":9,\"y\":15},{\"x\":13,\"y\":15},{\"x\":18,\"y\":15},{\"x\":0,\"y\":14},{\"x\":3,\"y\":14},{\"x\":5,\"y\":14},{\"x\":6,\"y\":14},{\"x\":7,\"y\":14},{\"x\":9,\"y\":14},{\"x\":11,\"y\":14},{\"x\":12,\"y\":14},{\"x\":13,\"y\":14},{\"x\":15,\"y\":14},{\"x\":18,\"y\":14},{\"x\":0,\"y\":13},{\"x\":3,\"y\":13},{\"x\":5,\"y\":13},{\"x\":13,\"y\":13},{\"x\":15,\"y\":13},{\"x\":18,\"y\":13},{\"x\":0,\"y\":12},{\"x\":1,\"y\":12},{\"x\":2,\"y\":12},{\"x\":3,\"y\":12},{\"x\":5,\"y\":12},{\"x\":7,\"y\":12},{\"x\":9,\"y\":12},{\"x\":11,\"y\":12},{\"x\":13,\"y\":12},{\"x\":15,\"y\":12},{\"x\":16,\"y\":12},{\"x\":17,\"y\":12},{\"x\":18,\"y\":12},{\"x\":7,\"y\":11},{\"x\":11,\"y\":11},{\"x\":0,\"y\":10},{\"x\":1,\"y\":10},{\"x\":2,\"y\":10},{\"x\":3,\"y\":10},{\"x\":5,\"y\":10},{\"x\":7,\"y\":10},{\"x\":9,\"y\":10},{\"x\":11,\"y\":10},{\"x\":13,\"y\":10},{\"x\":15,\"y\":10},{\"x\":16,\"y\":10},{\"x\":17,\"y\":10},{\"x\":18,\"y\":10},{\"x\":0,\"y\":9},{\"x\":3,\"y\":9},{\"x\":5,\"y\":9},{\"x\":13,\"y\":9},{\"x\":15,\"y\":9},{\"x\":18,\"y\":9},{\"x\":0,\"y\":8},{\"x\":3,\"y\":8},{\"x\":5,\"y\":8},{\"x\":7,\"y\":8},{\"x\":8,\"y\":8},{\"x\":9,\"y\":8},{\"x\":10,\"y\":8},{\"x\":11,\"y\":8},{\"x\":13,\"y\":8},{\"x\":15,\"y\":8},{\"x\":18,\"y\":8},{\"x\":0,\"y\":7},{\"x\":9,\"y\":7},{\"x\":18,\"y\":7},{\"x\":0,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":7,\"y\":6},{\"x\":9,\"y\":6},{\"x\":11,\"y\":6},{\"x\":12,\"y\":6},{\"x\":13,\"y\":6},{\"x\":15,\"y\":6},{\"x\":16,\"y\":6},{\"x\":18,\"y\":6},{\"x\":0,\"y\":5},{\"x\":3,\"y\":5},{\"x\":15,\"y\":5},{\"x\":18,\"y\":5},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":3,\"y\":4},{\"x\":5,\"y\":4},{\"x\":7,\"y\":4},{\"x\":8,\"y\":4},{\"x\":9,\"y\":4},{\"x\":10,\"y\":4},{\"x\":11,\"y\":4},{\"x\":13,\"y\":4},{\"x\":15,\"y\":4},{\"x\":17,\"y\":4},{\"x\":18,\"y\":4},{\"x\":0,\"y\":3},{\"x\":5,\"y\":3},{\"x\":9,\"y\":3},{\"x\":13,\"y\":3},{\"x\":18,\"y\":3},{\"x\":0,\"y\":2},{\"x\":2,\"y\":2},{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":6,\"y\":2},{\"x\":7,\"y\":2},{\"x\":9,\"y\":2},{\"x\":11,\"y\":2},{\"x\":12,\"y\":2},{\"x\":13,\"y\":2},{\"x\":14,\"y\":2},{\"x\":15,\"y\":2},{\"x\":16,\"y\":2},{\"x\":18,\"y\":2},{\"x\":0,\"y\":1},{\"x\":18,\"y\":1},{\"x\":0,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":6,\"y\":0},{\"x\":7,\"y\":0},{\"x\":8,\"y\":0},{\"x\":9,\"y\":0},{\"x\":10,\"y\":0},{\"x\":11,\"y\":0},{\"x\":12,\"y\":0},{\"x\":13,\"y\":0},{\"x\":14,\"y\":0},{\"x\":15,\"y\":0},{\"x\":16,\"y\":0},{\"x\":18,\"y\":0}]},\"you\":{\"id\":\"3e02eeed-67f2-4ded-8c7a-a554a5c8c588\",\"name\":\"local\",\"latency\":\"0\",\"health\":88,\"body\":[{\"x\":4,\"y\":19},{\"x\":4,\"y\":18},{\"x\":4,\"y\":17}],\"head\":{\"x\":4,\"y\":19},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}}")

	state := b.GameState{}
	_ = json.Unmarshal(g, &state)

	b := BuildBoard(state)

	if _, ok := b.Heads[MeId]; !ok {
		panic("Did not add me to board!")
	}
}

func TestCartesianProduct(t *testing.T) {
	//t.Skip()
	g := []byte("{\"game\":{\"id\":\"543b2935-7a1b-43fb-a6fb-42503ba7a28f\",\"ruleset\":{\"name\":\"wrapped\",\"version\":\"cli\",\"settings\":{\"foodSpawnChance\":15,\"minimumFood\":1,\"hazardDamagePerTurn\":100,\"hazardMap\":\"\",\"hazardMapAuthor\":\"\",\"royale\":{\"shrinkEveryNTurns\":25},\"squad\":{\"allowBodyCollisions\":false,\"sharedElimination\":false,\"sharedHealth\":false,\"sharedLength\":false}}},\"map\":\"arcade_maze\",\"timeout\":500,\"source\":\"\"},\"turn\":15,\"board\":{\"height\":21,\"width\":19,\"snakes\":[{\"id\":\"ca19ce3a-5199-4fdf-9405-7a0bec40dee0\",\"name\":\"tiam\",\"latency\":\"0\",\"health\":85,\"body\":[{\"x\":14,\"y\":12},{\"x\":14,\"y\":11},{\"x\":13,\"y\":11}],\"head\":{\"x\":14,\"y\":12},\"length\":3,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}},{\"id\":\"053a08e1-1814-495b-85ce-9dde707db720\",\"name\":\"local\",\"latency\":\"0\",\"health\":94,\"body\":[{\"x\":12,\"y\":8},{\"x\":12,\"y\":9},{\"x\":11,\"y\":9},{\"x\":10,\"y\":9}],\"head\":{\"x\":12,\"y\":8},\"length\":4,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}],\"food\":[{\"x\":14,\"y\":7}],\"hazards\":[{\"x\":0,\"y\":20},{\"x\":2,\"y\":20},{\"x\":3,\"y\":20},{\"x\":4,\"y\":20},{\"x\":5,\"y\":20},{\"x\":6,\"y\":20},{\"x\":7,\"y\":20},{\"x\":8,\"y\":20},{\"x\":9,\"y\":20},{\"x\":10,\"y\":20},{\"x\":11,\"y\":20},{\"x\":12,\"y\":20},{\"x\":13,\"y\":20},{\"x\":14,\"y\":20},{\"x\":15,\"y\":20},{\"x\":16,\"y\":20},{\"x\":18,\"y\":20},{\"x\":0,\"y\":19},{\"x\":9,\"y\":19},{\"x\":18,\"y\":19},{\"x\":0,\"y\":18},{\"x\":2,\"y\":18},{\"x\":3,\"y\":18},{\"x\":5,\"y\":18},{\"x\":6,\"y\":18},{\"x\":7,\"y\":18},{\"x\":9,\"y\":18},{\"x\":11,\"y\":18},{\"x\":12,\"y\":18},{\"x\":13,\"y\":18},{\"x\":15,\"y\":18},{\"x\":16,\"y\":18},{\"x\":18,\"y\":18},{\"x\":0,\"y\":17},{\"x\":18,\"y\":17},{\"x\":0,\"y\":16},{\"x\":2,\"y\":16},{\"x\":3,\"y\":16},{\"x\":5,\"y\":16},{\"x\":7,\"y\":16},{\"x\":8,\"y\":16},{\"x\":9,\"y\":16},{\"x\":10,\"y\":16},{\"x\":11,\"y\":16},{\"x\":13,\"y\":16},{\"x\":15,\"y\":16},{\"x\":16,\"y\":16},{\"x\":18,\"y\":16},{\"x\":0,\"y\":15},{\"x\":5,\"y\":15},{\"x\":9,\"y\":15},{\"x\":13,\"y\":15},{\"x\":18,\"y\":15},{\"x\":0,\"y\":14},{\"x\":3,\"y\":14},{\"x\":5,\"y\":14},{\"x\":6,\"y\":14},{\"x\":7,\"y\":14},{\"x\":9,\"y\":14},{\"x\":11,\"y\":14},{\"x\":12,\"y\":14},{\"x\":13,\"y\":14},{\"x\":15,\"y\":14},{\"x\":18,\"y\":14},{\"x\":0,\"y\":13},{\"x\":3,\"y\":13},{\"x\":5,\"y\":13},{\"x\":13,\"y\":13},{\"x\":15,\"y\":13},{\"x\":18,\"y\":13},{\"x\":0,\"y\":12},{\"x\":1,\"y\":12},{\"x\":2,\"y\":12},{\"x\":3,\"y\":12},{\"x\":5,\"y\":12},{\"x\":7,\"y\":12},{\"x\":9,\"y\":12},{\"x\":11,\"y\":12},{\"x\":13,\"y\":12},{\"x\":15,\"y\":12},{\"x\":16,\"y\":12},{\"x\":17,\"y\":12},{\"x\":18,\"y\":12},{\"x\":7,\"y\":11},{\"x\":11,\"y\":11},{\"x\":0,\"y\":10},{\"x\":1,\"y\":10},{\"x\":2,\"y\":10},{\"x\":3,\"y\":10},{\"x\":5,\"y\":10},{\"x\":7,\"y\":10},{\"x\":9,\"y\":10},{\"x\":11,\"y\":10},{\"x\":13,\"y\":10},{\"x\":15,\"y\":10},{\"x\":16,\"y\":10},{\"x\":17,\"y\":10},{\"x\":18,\"y\":10},{\"x\":0,\"y\":9},{\"x\":3,\"y\":9},{\"x\":5,\"y\":9},{\"x\":13,\"y\":9},{\"x\":15,\"y\":9},{\"x\":18,\"y\":9},{\"x\":0,\"y\":8},{\"x\":3,\"y\":8},{\"x\":5,\"y\":8},{\"x\":7,\"y\":8},{\"x\":8,\"y\":8},{\"x\":9,\"y\":8},{\"x\":10,\"y\":8},{\"x\":11,\"y\":8},{\"x\":13,\"y\":8},{\"x\":15,\"y\":8},{\"x\":18,\"y\":8},{\"x\":0,\"y\":7},{\"x\":9,\"y\":7},{\"x\":18,\"y\":7},{\"x\":0,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":7,\"y\":6},{\"x\":9,\"y\":6},{\"x\":11,\"y\":6},{\"x\":12,\"y\":6},{\"x\":13,\"y\":6},{\"x\":15,\"y\":6},{\"x\":16,\"y\":6},{\"x\":18,\"y\":6},{\"x\":0,\"y\":5},{\"x\":3,\"y\":5},{\"x\":15,\"y\":5},{\"x\":18,\"y\":5},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":3,\"y\":4},{\"x\":5,\"y\":4},{\"x\":7,\"y\":4},{\"x\":8,\"y\":4},{\"x\":9,\"y\":4},{\"x\":10,\"y\":4},{\"x\":11,\"y\":4},{\"x\":13,\"y\":4},{\"x\":15,\"y\":4},{\"x\":17,\"y\":4},{\"x\":18,\"y\":4},{\"x\":0,\"y\":3},{\"x\":5,\"y\":3},{\"x\":9,\"y\":3},{\"x\":13,\"y\":3},{\"x\":18,\"y\":3},{\"x\":0,\"y\":2},{\"x\":2,\"y\":2},{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":6,\"y\":2},{\"x\":7,\"y\":2},{\"x\":9,\"y\":2},{\"x\":11,\"y\":2},{\"x\":12,\"y\":2},{\"x\":13,\"y\":2},{\"x\":14,\"y\":2},{\"x\":15,\"y\":2},{\"x\":16,\"y\":2},{\"x\":18,\"y\":2},{\"x\":0,\"y\":1},{\"x\":18,\"y\":1},{\"x\":0,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":6,\"y\":0},{\"x\":7,\"y\":0},{\"x\":8,\"y\":0},{\"x\":9,\"y\":0},{\"x\":10,\"y\":0},{\"x\":11,\"y\":0},{\"x\":12,\"y\":0},{\"x\":13,\"y\":0},{\"x\":14,\"y\":0},{\"x\":15,\"y\":0},{\"x\":16,\"y\":0},{\"x\":18,\"y\":0}]},\"you\":{\"id\":\"053a08e1-1814-495b-85ce-9dde707db720\",\"name\":\"local\",\"latency\":\"0\",\"health\":94,\"body\":[{\"x\":12,\"y\":8},{\"x\":12,\"y\":9},{\"x\":11,\"y\":9},{\"x\":10,\"y\":9}],\"head\":{\"x\":12,\"y\":8},\"length\":4,\"shout\":\"\",\"squad\":\"\",\"customizations\":{\"color\":\"#002080\",\"head\":\"evil\",\"tail\":\"fat-rattle\"}}}")

	state := b.GameState{}
	_ = json.Unmarshal(g, &state)

	b := BuildBoard(state)
	states := GetCartesianProductOfMoves(&b)
	if len(states) != 1 {
		panic("Did not get correct states!")
	}
}
