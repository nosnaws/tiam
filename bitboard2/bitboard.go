package bitboard

import (
	num "github.com/shabbyrobe/go-num"

	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/nosnaws/tiam/moveset"
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

type SnakeMoveSet struct {
	Id  string
	Set moveset.MoveSet
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
	isRoyale        bool
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
		isRoyale:        state.Game.Map == "royale",
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

// BOARD
// 111 | 112 | 113 | 114 | 115 | 116 | 117 | 118 | 119 | 120 | 121
// 100 | 101 | 102 | 103 | 104 | 105 | 106 | 107 | 108 | 109 | 110
// 089 | 090 | 091 | 092 | 093 | 094 | 095 | 096 | 097 | 098 | 099
// 078 | 079 | 080 | 081 | 082 | 083 | 084 | 085 | 086 | 087 | 088
// 067 | 068 | 069 | 070 | 071 | 072 | 073 | 074 | 075 | 076 | 077
// 055 | 056 | 057 | 058 | 059 | 061 | 062 | 063 | 064 | 065 | 066
// 044 | 045 | 046 | 047 | 048 | 049 | 050 | 051 | 052 | 053 | 054
// 033 | 034 | 035 | 036 | 037 | 038 | 039 | 040 | 041 | 042 | 043
// 022 | 023 | 024 | 025 | 026 | 027 | 028 | 029 | 030 | 031 | 032
// 011 | 012 | 013 | 014 | 015 | 016 | 017 | 018 | 019 | 020 | 021
// 000 | 001 | 002 | 003 | 004 | 005 | 006 | 007 | 008 | 009 | 010

var BOTTOM_MASK = num.U128FromRaw(
	0b0000000000000000000000000000000000000000000000000000000000000000, // HI
	0b0000000000000000000000000000000000000000000000000000011111111111, // LO
)

// This is the 11 most sig bits + 7 unused bits
var TOP_MASK = num.U128FromRaw(
	0b1111111111111111110000000000000000000000000000000000000000000000, // HI
	0b0000000000000000000000000000000000000000000000000000000000000000, // LO
)

var RIGHT_MASK = num.U128FromRaw(
	0b0000000100000000001000000000010000000000100000000001000000000010, // HI
	0b0000000001000000000010000000000100000000001000000000010000000000, // LO
)

var LEFT_MASK = num.U128FromRaw(
	0b0000000000000000010000000000100000000001000000000010000000000100, // HI
	0b0000000010000000000100000000001000000000010000000000100000000001, // LO
)

func (bb *BitBoard) IsSnakeOnWall(snakeId string) bool {
	if bb.isWrapped {
		return false
	}

	if snake := bb.GetSnake(snakeId); snake != nil {
		head := snake.getHeadBoard()

		if head.And(LEFT_MASK).BitLen() > 0 {
			return true
		}

		if head.And(TOP_MASK).BitLen() > 0 {
			return true
		}

		if head.And(RIGHT_MASK).BitLen() > 0 {
			return true
		}

		if head.And(BOTTOM_MASK).BitLen() > 0 {
			return true
		}
	}
	return false
}

func (bb *BitBoard) GetMoves(snakeId string) SnakeMoveSet {
	width := uint(bb.width)
	moves := SnakeMoveSet{
		Id: snakeId,
	}
	ms := moveset.Create()

	snake := bb.GetSnake(snakeId)
	if !snake.IsAlive() {
		return moves
	}

	//headIndex := snake.GetHeadIndex()
	headBoard := snake.getHeadBoard()
	//fmt.Println("HEAD BOARD")
	//bb.printBoard(headBoard)

	leftBoard := headBoard.Rsh(1)
	//fmt.Println("LEFT BOARD")
	//bb.printBoard(leftBoard)
	rightBoard := headBoard.Lsh(1)
	//fmt.Println("RIGHT BOARD")
	//bb.printBoard(rightBoard)
	upBoard := headBoard.Lsh(width)
	//fmt.Println("UP BOARD")
	//bb.printBoard(upBoard)
	downBoard := headBoard.Rsh(width)
	//fmt.Println("DOWN BOARD")
	//bb.printBoard(downBoard)

	// LEFT
	if headBoard.And(LEFT_MASK).BitLen() == 0 {
		if leftBoard.And(bb.empty).BitLen() > 0 {
			ms = moveset.SetLeft(ms)
		}
	}

	// RIGHT
	if headBoard.And(RIGHT_MASK).BitLen() == 0 {
		if rightBoard.And(bb.empty).BitLen() > 0 {
			ms = moveset.SetRight(ms)
		}
	}

	// UP
	if headBoard.And(TOP_MASK).BitLen() == 0 {
		if upBoard.And(bb.empty).BitLen() > 0 {
			ms = moveset.SetUp(ms)
		}
	}

	// DOWN
	if headBoard.And(BOTTOM_MASK).BitLen() == 0 {
		if downBoard.And(bb.empty).BitLen() > 0 {
			ms = moveset.SetDown(ms)
		}
	}

	if moveset.IsEmpty(ms) {
		ms = moveset.SetLeft(ms)
	}

	moves.Set = ms

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
	snake.moveHead(newHead, dir, uint(bb.width))
}

func (bb *BitBoard) isSnakeOutOfBounds(id string, headIdx int, ms moveset.MoveSet) bool {
	if isDirOutOfBounds(MoveSetToDir(ms), headIdx, bb.width, bb.height, bb.isWrapped) {
		return true
	}

	return false
}

func (bb *BitBoard) AdvanceTurn(allMoves []SnakeMoveSet) {
	deadSnakes := []string{}
	bb.turn += 1

	for _, move := range allMoves {
		snake := bb.GetSnake(move.Id)
		if snake == nil {
			continue
		}

		// moved out of bounds, need to handle this here before state gets messed up
		if bb.isSnakeOutOfBounds(move.Id, snake.GetHeadIndex(), move.Set) {
			bb.killSnake(move.Id)
			continue
		}

		//fmt.Println("MOVING", move.Id)
		//bb.printBoard(snake.headBoard)
		bb.moveSnake(move.Id, MoveSetToDir(move.Set))
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
	for _, move := range allMoves {
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

	if bb.isRoyale {
		bb.SpawnHazardsRoyale()
	}
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
		isRoyale:     bb.isRoyale,
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
		//if !snake.IsAlive() {
		//continue
		//}

		if snake.board.Xor(board.GetSnake(id).board).BitLen() > 0 {
			return false
		}
	}

	return true
}

// returns the Dir version ("left", "right", "up", "down")
// of the moveset. If more than one is set, it returns the first found
// in the order displayed above.
func MoveSetToDir(ms moveset.MoveSet) Dir {
	if moveset.HasLeft(ms) {
		return Left
	}
	if moveset.HasRight(ms) {
		return Right
	}
	if moveset.HasUp(ms) {
		return Up
	}
	if moveset.HasDown(ms) {
		return Down
	}
	return ""
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
