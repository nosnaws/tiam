package bitboard

import (
	"math/big"

	api "github.com/nosnaws/tiam/battlesnake"
)

const MeId = 0

// bitboard
// use an array of uint64s to represent different parts of the game state
// an array is needed because the standard 11x11 board does not fit into 64 bits
// other option is big.Int
//
// board per snake, keep track of head and tail (use a bit to denote that a snake ate in the last turn?)
// board for food
// board for hazards
// empty board
// tail movement can be tracked in the empty board, the tile won't be empty if the snake ate the turn before
//

type SnakeMove struct {
	id  int
	dir Dir
}

type BitBoard struct {
	food         *big.Int
	hazards      *big.Int
	snakes       []*snake
	width        int
	height       int
	empty        *big.Int
	hazardDamage int
	isWrapped    bool
	turn         int
}

func (bb *BitBoard) createEmptyBoard() *big.Int {
	emptyBoard := big.NewInt(0)
	emptyBoard.Not(emptyBoard)

	for _, snake := range bb.snakes {
		emptyBoard.Xor(emptyBoard, snake.board)
	}

	// TODO: make hazards valid moves
	emptyBoard.Xor(emptyBoard, bb.hazards)

	return emptyBoard
}

func CreateBitBoard(state api.GameState) BitBoard {
	me := createSnake(state.You, state.Board.Width)
	snakes := []*snake{me}

	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		snakes = append(snakes, createSnake(snake, state.Board.Width))
	}

	foodBoard := big.NewInt(0)
	for _, p := range state.Board.Food {
		foodBoard.SetBit(foodBoard, getIndex(p, state.Board.Width), 1)
	}

	hazardsBoard := big.NewInt(0)
	for _, p := range state.Board.Hazards {
		hazardsBoard.SetBit(hazardsBoard, getIndex(p, state.Board.Width), 1)
	}

	bb := BitBoard{
		food:         foodBoard,
		hazards:      hazardsBoard,
		snakes:       snakes,
		width:        state.Board.Width,
		height:       state.Board.Height,
		hazardDamage: int(state.Game.Ruleset.Settings.HazardDamagePerTurn),
		isWrapped:    state.Game.Ruleset.Name == "wrapped",
		turn:         state.Turn,
	}
	bb.empty = bb.createEmptyBoard()

	return bb
}

func (bb *BitBoard) getSnake(id int) *snake {
	return bb.snakes[id]
}

func (bb *BitBoard) moveSnake(id int, dir Dir) {
	snake := bb.getSnake(id)
	head := snake.getHeadIndex()

	newHead := indexInDirection(dir, head, bb.width, bb.height, bb.isWrapped)

	// move tail
	snake.moveTail()

	// move head
	snake.moveHead(newHead)
}

func (bb *BitBoard) advanceTurn(moves []SnakeMove) {
	for _, move := range moves {
		snake := bb.getSnake(move.id)

		// moved out of bounds, need to handle this here before state gets messed up
		if isDirOutOfBounds(move.dir, snake.getHeadIndex(), bb.width, bb.height, bb.isWrapped) {
			snake.kill()
			continue
		}

		bb.moveSnake(move.id, move.dir)
		bb.getSnake(move.id).health -= 1
	}

	// TODO: hazard damage

	// feed snakes
	for _, move := range moves {
		snake := bb.getSnake(move.id)
		if !snake.isAlive() {
			continue
		}

		foodCopy := big.NewInt(0).Set(bb.food)
		foundFood := foodCopy.And(foodCopy, snake.board).BitLen() > 0

		if foundFood {
			snake.feed()

			// remove food from board
			bb.food.SetBit(bb.food, snake.getHeadIndex(), 0)
		}
	}

	// TODO: spawn food?

	// kill snakes
	for _, snake := range bb.snakes {
		if !snake.isAlive() {
			continue
		}

		if snake.health < 1 {
			snake.kill()
			continue
		}

		// collided with itself

		// collided with another snake

		// head to head loss

	}
}
