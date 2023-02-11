package bitboard

import (
	"fmt"
	"testing"

	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/shabbyrobe/go-num"
)

func TestAdvanceBoardTurnDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ _
	// s s h
	me := api.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}

	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	if board.GetSnake(me.ID).health != 99 {
		fmt.Println(board.GetSnake(me.ID).health)
		board.Print()
		panic("Did not decrement health properly!")
	}
}

func TestHeadToHead(t *testing.T) {
	//t.Skip()
	// _ _ _ _ _
	// _ s s s h
	// s s s e f
	// _ _ _ _ _
	// _ _ _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 4, Y: 3}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}},
	}
	two := api.Battlesnake{
		ID:     "two",
		Health: 100,
		Body:   []api.Coord{{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}},
	}
	state := api.GameState{
		Board: api.Board{
			Snakes: []api.Battlesnake{me, two},
			Height: 5,
			Width:  5,
			Food:   []api.Coord{{X: 4, Y: 2}},
		},
		You: me,
	}
	board := CreateBitBoard(state)

	moves := []SnakeMove{
		{Id: me.ID, Dir: Down},
		{Id: two.ID, Dir: Right},
	}
	board.AdvanceTurn(moves)

	if !board.IsGameOver() {
		board.Print()
		panic("game did not end!")
	}
}

func TestAdvanceBoardHazardDamage(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ x
	// s s h
	me := api.Battlesnake{
		ID:     "me",
		Health: 50,
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes:  []api.Battlesnake{me},
			Height:  3,
			Width:   3,
			Hazards: []api.Coord{{X: 2, Y: 1}},
		},
		You: me,
		Game: api.Game{
			Ruleset: api.Ruleset{
				Settings: api.Settings{
					HazardDamagePerTurn: 100,
				},
			},
		},
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	if board.GetSnake(me.ID) != nil {
		fmt.Println(board.GetSnake(me.ID).health)
		board.Print()
		panic("Snake did not die!")
	}
}

func TestAdvanceBoardEatFood(t *testing.T) {
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
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	snake := board.GetSnake(me.ID)
	oldTail := snake.getTailIndex()
	if snake.health != 100 {
		board.Print()
		panic("Snake did not eat!")
	}
	if snake.Length != 4 {
		panic("Snake did not grow!")
	}

	if snake.GetHeadIndex() != 5 {
		board.Print()
		panic("swapped head and tail!")
	}

	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	snake = board.GetSnake(me.ID)
	if snake.health != 99 {
		panic("Did not reduce health!")
	}
	if snake.getTailIndex() != oldTail {
		panic("Moved tail!")
	}
	if snake.GetHeadIndex() != 8 {
		board.Print()
		panic("did not move head properly!")
	}
}

func TestAdvanceBoardOutOfBounds(t *testing.T) {
	//t.Skip()
	// right side
	// _ _ _
	// _ _ _
	// s s h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Right}})

	if board.GetSnake(me.ID) != nil {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// left side
	// _ _ _
	// _ _ _
	// h s s
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 0, Y: 0},
		Body:   []api.Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
	}
	state = api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = CreateBitBoard(state)

	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Left}})

	if board.GetSnake(me.ID) != nil {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// top side
	// h _ _
	// s _ _
	// s _ _
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 0, Y: 2},
		Body:   []api.Coord{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
	}
	state = api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	if board.GetSnake(me.ID) != nil {
		panic("Did not remove snake from board!")
	}

	//t.Skip()
	// bottom side
	// _ _ s
	// _ _ s
	// _ _ h
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Head:   api.Coord{X: 2, Y: 0},
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
	}
	state = api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Down}})

	if board.GetSnake(me.ID) != nil {
		panic("Did not remove snake from board!")
	}
}

func TestAdvanceBoardSnakeCollision(t *testing.T) {
	//t.Skip()
	// o o _
	// o e _
	// s s h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	enemy := api.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []api.Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
	}
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}, {Id: enemy.ID, Dir: Down}})

	if board.GetSnake(enemy.ID) != nil {
		fmt.Println(board)
		panic("Should have removed enemy from the board!")
	}

	if board.GetSnake(me.ID) == nil {
		panic("Should not have removed me!")
	}
}

func TestAdvanceBoardSelfCollision(t *testing.T) {
	//t.Skip()
	// _ s _
	// _ s h
	// _ s s
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Left}})

	if board.GetSnake(me.ID) != nil {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardFollowTail(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ s s
	// _ s h
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}},
	}
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 2, Y: 1}},
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	if board.GetSnake(me.ID) == nil {
		panic("Snake was killed!")
	}

	if board.GetSnake(me.ID).GetHeadIndex() != 5 {
		board.Print()
		panic("Did not move snake!")
	}

	// Follow other snake tail
	// _ e e
	// _ s e
	// _ s h
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	}
	enemy := api.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []api.Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}},
	}
	state = api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}, {Id: enemy.ID, Dir: Left}})

	if board.GetSnake(me.ID) == nil {
		panic("Snake was killed!")
	}

	if board.GetSnake(me.ID).GetHeadIndex() != 5 {
		panic("Did not move snake!")
	}

	// Follow other snake tail
	// f e e
	// _ s e
	// _ s h
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	}
	enemy = api.Battlesnake{
		ID:     "enemy",
		Health: 100,
		Body:   []api.Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 1}},
	}
	state = api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, enemy},
			Height: 3,
			Width:  3,
			Food:   []api.Coord{{X: 0, Y: 2}},
		},
		You: me,
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}, {Id: enemy.ID, Dir: Left}})

	if board.GetSnake(me.ID) != nil {
		panic("Snake was not killed!")
	}
}

func TestAdvanceBoardMoveOnNeck(t *testing.T) {
	//t.Skip()
	// _ s _
	// _ s h
	// _ s s
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state := api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Down}})

	if board.GetSnake(me.ID) != nil {
		panic("Should have killed me!")
	}
}

func TestAdvanceBoardWrapped(t *testing.T) {
	//t.Skip()
	// _ _ _
	// _ _ h
	// _ s s
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}},
	}
	state := api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Right}})
	headBoard := board.GetSnake(me.ID).headBoard

	if board.GetSnake(me.ID) == nil {
		panic("Should not have killed me!")
	}

	if board.GetSnake(me.ID).GetHeadIndex() != 3 {
		panic("Did not wrapped to right!")
	}

	testBoard := num.U128From16(0)
	testBoard = testBoard.SetBit(3, 1)
	if testBoard.And(headBoard).BitLen() == 0 {
		board.printBoard(headBoard)
		panic("Did not update headboard")
	}

	// _ _ h
	// _ _ s
	// _ _ s
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}},
	}
	state = api.GameState{
		Turn: 3,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Up}})

	if board.GetSnake(me.ID) == nil {
		panic("Should not have killed me!")
	}

	if board.GetSnake(me.ID).GetHeadIndex() != 2 {
		panic("Did not wrapped to bottom!")
	}

	// _ s _
	// _ s _
	// _ h _
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
	}
	state = api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Down}})

	if board.GetSnake(me.ID) == nil {
		panic("Should not have killed me!")
	}

	if board.GetSnake(me.ID).GetHeadIndex() != 7 {
		panic("Did not wrapped to top!")
	}

	// _ _ _
	// s s h
	// _ _ _
	me = api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 1}},
	}
	state = api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board = CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{{Id: me.ID, Dir: Right}})

	if board.GetSnake(me.ID) == nil {
		panic("Should not have killed me!")
	}

	if board.GetSnake(me.ID).GetHeadIndex() != 3 {
		panic("Did not wrapped to left!")
	}
}

func TestGetMoves(t *testing.T) {
	// _ _ _
	// s s h
	// _ _ _
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 1}},
	}
	state := api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me},
			Height: 3,
			Width:  3,
		},
		You: me,
		Game: api.Game{
			Ruleset: api.Ruleset{
				Name: "wrapped",
			},
		},
	}
	board := CreateBitBoard(state)
	moves := board.GetMoves(me.ID)

	if len(moves) != 3 {
		board.Print()
		fmt.Println(moves)
		panic("wrong number of moves!")
	}
}

func TestHeadToHeadOnFood(t *testing.T) {
	// s _ _
	// h f e
	// _ _ s
	me := api.Battlesnake{
		ID:     "me",
		Health: 100,
		Body:   []api.Coord{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 2}},
	}
	s2 := api.Battlesnake{
		ID:     "s2",
		Health: 100,
		Body:   []api.Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 2, Y: 0}},
	}
	state := api.GameState{
		Turn: 4,
		Board: api.Board{
			Snakes: []api.Battlesnake{me, s2},
			Height: 3,
			Width:  3,
			Food: []api.Coord{
				{X: 1, Y: 1},
			},
		},
		You: me,
	}
	board := CreateBitBoard(state)
	board.AdvanceTurn([]SnakeMove{
		{Id: me.ID, Dir: Right},
		{Id: s2.ID, Dir: Left},
	})

	if !board.IsGameOver() {
		panic("did not end game!")
	}

	if board.GetSnake(me.ID) != nil {
		panic("my snake should not be alive")
	}

	if _, ok := board.Snakes[me.ID]; ok {
		panic("my snake should not be alive")
	}

	if board.GetSnake(s2.ID) != nil {
		panic("other snake should not be alive")
	}

	if board.food.AsUint64() == 0 {
		fmt.Println(board.food)
		board.printBoard(board.food)
		panic("food was removed!")
	}

}
