package bitboard

import (
	num "github.com/shabbyrobe/go-num"

	api "github.com/nosnaws/tiam/battlesnake"
)

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
	Id  string
	Dir Dir
}

type BitBoard struct {
	food    num.U128
	hazards num.U128
	// I need to makes Snakes a map so that I can remove snakes when they die.
	// Currently, removing a snake from the array could cause the ids of other snakes to change, since it's an array.
	// By converting to a map, I can just use the ID from the API
	// This way, when a snake dies and the API does not send that snake in the request, the state will match up with what I have in my tree since it is also being removed
	Snakes          map[string]*snake
	meId            string
	width           int
	height          int
	empty           num.U128
	hazardDamage    int
	hazardSpawnTime int
	foodSpawnChance int
	minimumFood     int
	totalFood       int
	isWrapped       bool
	turn            int
}

func (bb *BitBoard) createEmptyBoard() num.U128 {
	emptyBoard := num.U128From16(0).Not()

	for _, snake := range bb.Snakes {
		if !snake.IsAlive() {
			continue
		}
		emptyBoard = emptyBoard.Xor(snake.board)

		if !snake.stackedTail() {
			emptyBoard = emptyBoard.SetBit(snake.getTailIndex(), 1)
		}
	}

	return emptyBoard
}

func CreateBitBoard(state api.GameState) *BitBoard {
	me := createSnake(state.You, state.Board.Width)
	snakes := make(map[string]*snake)
	snakes[state.You.ID] = me

	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		snakes[snake.ID] = createSnake(snake, state.Board.Width)
	}

	foodBoard := createBoard(state.Board.Food, state.Board.Width)
	hazardsBoard := createBoard(state.Board.Hazards, state.Board.Width)

	bb := BitBoard{
		food:            foodBoard,
		hazards:         hazardsBoard,
		Snakes:          snakes,
		meId:            state.You.ID,
		width:           state.Board.Width,
		height:          state.Board.Height,
		hazardDamage:    int(state.Game.Ruleset.Settings.HazardDamagePerTurn),
		hazardSpawnTime: int(state.Game.Ruleset.Settings.Royale.ShrinkEveryNTurns),
		minimumFood:     int(state.Game.Ruleset.Settings.MinimumFood),
		foodSpawnChance: int(state.Game.Ruleset.Settings.FoodSpawnChance),
		totalFood:       len(state.Board.Food),
		isWrapped:       state.Game.Ruleset.Name == "wrapped",
		turn:            state.Turn,
	}
	bb.empty = bb.createEmptyBoard()

	//bb.printBoard(snakes["me"].headBoard)

	return &bb
}

func createBoard(coords []api.Coord, width int) num.U128 {
	board := num.U128From16(0)
	for _, p := range coords {
		board = board.SetBit(getIndex(p, width), 1)
	}

	return board
}

func (bb *BitBoard) IsGameOver() bool {
	return len(bb.Snakes) < 2
	//numSnakesAlive := 0
	//for _, snake := range bb.Snakes {
	//if snake.IsAlive() {
	//numSnakesAlive += 1
	//}
	//}

	//if numSnakesAlive > 1 {
	//return false
	//}

	//return true
}

func (bb *BitBoard) GetMoves(snakeId string) []SnakeMove {
	moves := []SnakeMove{}
	snake := bb.GetSnake(snakeId)
	if !snake.IsAlive() {
		return moves
	}

	headIndex := snake.GetHeadIndex()

	if !isDirOutOfBounds(Left, headIndex, bb.width, bb.height, bb.isWrapped) {
		leftBoard := num.U128From16(0)
		leftIndex := indexInDirection(Left, headIndex, bb.width, bb.height, bb.isWrapped)
		leftBoard = leftBoard.SetBit(leftIndex, 1)

		if leftBoard.And(bb.empty).BitLen() > 0 {
			moves = append(moves, SnakeMove{Id: snakeId, Dir: Left})
		}
	}

	if !isDirOutOfBounds(Right, headIndex, bb.width, bb.height, bb.isWrapped) {
		rightBoard := num.U128From16(0)
		rightIndex := indexInDirection(Right, headIndex, bb.width, bb.height, bb.isWrapped)
		rightBoard = rightBoard.SetBit(rightIndex, 1)

		if rightBoard.And(bb.empty).BitLen() > 0 {
			moves = append(moves, SnakeMove{Id: snakeId, Dir: Right})
		}
	}

	if !isDirOutOfBounds(Up, headIndex, bb.width, bb.height, bb.isWrapped) {
		upBoard := num.U128From16(0)
		upIndex := indexInDirection(Up, headIndex, bb.width, bb.height, bb.isWrapped)
		upBoard = upBoard.SetBit(upIndex, 1)

		if upBoard.And(bb.empty).BitLen() > 0 {
			moves = append(moves, SnakeMove{Id: snakeId, Dir: Up})
		}
	}

	if !isDirOutOfBounds(Down, headIndex, bb.width, bb.height, bb.isWrapped) {
		downBoard := num.U128From16(0)
		downIndex := indexInDirection(Down, headIndex, bb.width, bb.height, bb.isWrapped)
		downBoard = downBoard.SetBit(downIndex, 1)

		if downBoard.And(bb.empty).BitLen() > 0 {
			moves = append(moves, SnakeMove{Id: snakeId, Dir: Down})
		}
	}

	if len(moves) < 1 {
		moves = append(moves, SnakeMove{Id: snakeId, Dir: Left})
	}

	return moves
}

func (bb *BitBoard) GetSnake(id string) *snake {
	if snake, ok := bb.Snakes[id]; ok {
		return snake
	}
	return nil
}

func (bb *BitBoard) moveSnake(id string, dir Dir) {
	snake := bb.GetSnake(id)
	head := snake.GetHeadIndex()

	newHead := indexInDirection(dir, head, bb.width, bb.height, bb.isWrapped)

	// move tail
	snake.moveTail()

	// move head
	snake.moveHead(newHead)
}

func (bb *BitBoard) AdvanceTurn(moves []SnakeMove) {
	deadSnakes := []string{}
	bb.turn += 1
	for _, move := range moves {
		snake := bb.GetSnake(move.Id)
		if snake == nil {
			continue
		}

		// moved out of bounds, need to handle this here before state gets messed up
		if isDirOutOfBounds(move.Dir, snake.GetHeadIndex(), bb.width, bb.height, bb.isWrapped) {
			//fmt.Println("OUT OF BOUND")
			bb.killSnake(move.Id)
			continue
		}

		//fmt.Println("MOVING", move.Id)
		//bb.printBoard(snake.headBoard)
		bb.moveSnake(move.Id, move.Dir)
		//fmt.Println("AFTER MVOING")
		//bb.printBoard(snake.headBoard)
		bb.GetSnake(move.Id).health -= 1
	}

	// kill snakes
	for id, snake := range bb.Snakes {
		if !snake.IsAlive() {
			continue
		}

		if snake.health < 1 {
			//bb.killSnake(id)
			//fmt.Println("OUT OF HEALTH")
			deadSnakes = append(deadSnakes, id)
			continue
		}

		// collision
		headBoard := snake.getHeadBoard()

		// i think this covers all collisions except head to heads
		if headBoard.And(bb.empty).BitLen() == 0 {
			//bb.killSnake(id)
			deadSnakes = append(deadSnakes, id)
			//fmt.Println("COLLISION")
			//bb.printBoard(headBoard)
			//bb.printBoard(bb.empty)

			// collided with itself

			// collided with another snake

			// head to head loss
		}

		for otherSId, otherS := range bb.Snakes {
			if !otherS.IsAlive() {
				continue
			}

			if id == otherSId {
				continue
			}

			otherHeadBoard := otherS.getHeadBoard()
			// head to head
			if otherHeadBoard.And(headBoard).BitLen() > 0 {
				if otherS.Length > snake.Length {
					//bb.killSnake(id)
					//fmt.Println("HEAD TO HEAD LOSE")
					deadSnakes = append(deadSnakes, id)
					continue
				}

				if otherS.Length == snake.Length {
					//fmt.Println("HEAD TO HEAD DRAW")
					deadSnakes = append(deadSnakes, id)
					deadSnakes = append(deadSnakes, otherSId)
					//bb.killSnake(id)
					//bb.killSnake(otherSId)
				}
			}
		}
	}

	for _, id := range deadSnakes {
		bb.killSnake(id)
	}

	// feed snakes
	for _, move := range moves {
		snake := bb.GetSnake(move.Id)
		if snake == nil {
			continue
		}

		//foodCopy := big.NewInt(0).Set(bb.food)
		foundFood := bb.food.And(snake.board).BitLen() > 0

		if foundFood {
			snake.feed()

			// remove food from board
			bb.food = bb.food.SetBit(snake.GetHeadIndex(), 0)
			bb.totalFood -= 1
		} else {
			// hazard damage is applied if food is not found
			headBoard := snake.getHeadBoard()
			if headBoard.And(bb.hazards).BitLen() > 0 {
				snake.health -= bb.hazardDamage
			}

			if snake.health <= 0 {
				//fmt.Println("HAZARD KILL")
				bb.killSnake(move.Id)
			}
		}
	}

	bb.SpawnHazardsRoyale()
	bb.SpawnFood()

	bb.empty = bb.createEmptyBoard()
}

func (bb *BitBoard) killSnake(id string) {
	snake := bb.GetSnake(id)
	if snake != nil {
		snake.kill()
		delete(bb.Snakes, id)
	}
}

func (bb *BitBoard) IsIndexOccupied(i int) bool {
	tester := num.U128From16(0)
	tester = tester.SetBit(i, 1)

	return tester.And(bb.empty).BitLen() == 0
}

func (bb *BitBoard) IsIndexFood(i int) bool {
	tester := num.U128From16(0)
	tester = tester.SetBit(i, 1)

	return tester.And(bb.food).BitLen() > 0
}

func (bb *BitBoard) IsIndexHazard(i int) bool {
	tester := num.U128From16(0)
	tester = tester.SetBit(i, 1)

	return tester.And(bb.hazards).BitLen() > 0
}

func (bb *BitBoard) Clone() *BitBoard {
	snakes := make(map[string]*snake, len(bb.Snakes))
	for id, snake := range bb.Snakes {
		snakes[id] = snake.clone()
	}

	food := bb.food
	hazards := bb.hazards
	empty := bb.empty

	return &BitBoard{
		food:         food,
		hazards:      hazards,
		empty:        empty,
		Snakes:       snakes,
		width:        bb.width,
		height:       bb.height,
		hazardDamage: bb.hazardDamage,
		isWrapped:    bb.isWrapped,
		turn:         bb.turn,
	}
}

func (bb *BitBoard) IsEqual(board *BitBoard) bool {
	if bb.empty.Xor(board.empty).BitLen() > 0 {
		return false
	}

	if len(bb.Snakes) != len(board.Snakes) {
		return false
	}

	for id, snake := range bb.Snakes {
		if !snake.IsAlive() {
			continue
		}

		if snake.board.Xor(board.GetSnake(id).board).BitLen() > 0 {
			return false
		}
	}

	return true
}

//func (bb *BitBoard) AdvanceWithExternal(state api.GameState) {
//newSnakes := make(map[string]api.Battlesnake)
//for _, s := range state.Board.Snakes {
//newSnakes[s.ID] = s
//}

//actualMoves := []SnakeMove{}

//for id, snake := range bb.Snakes {
//if nSnake, ok := newSnakes[id]; ok {
//actualMove := SnakeMove{
//Id:  id,
//Dir: bb.GetLastSnakeMoveFromExternal(nSnake),
//}
//actualMoves = append(actualMoves, actualMove)
//} else {
//// remove snakes that died
//// i think this won't mess up turn resolution...
//bb.GetSnake(bId).kill()
//}
//}

//bb.AdvanceTurn(actualMoves)
//bb.food = createBoard(state.Board.Food, bb.width)
//bb.hazards = createBoard(state.Board.Hazards, bb.width)
//}
