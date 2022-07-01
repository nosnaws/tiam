package game

import (
	"fmt"
	"math/rand"
)

const MeId = SnakeId(1)

type FastBoard struct {
	list      []Tile
	ids       map[string]SnakeId
	Heads     map[SnakeId]uint16
	Lengths   map[SnakeId]uint8
	Healths   map[SnakeId]int8
	width     uint16
	height    uint16
	isWrapped bool
}

type SnakeMove struct {
	Id  SnakeId
	Dir Move
}

type idToNewHead struct {
	id    SnakeId
	index uint16
}

func BuildBoard(state GameState) FastBoard {
	height := uint16(state.Board.Height)
	width := uint16(state.Board.Width)
	listLength := height * width

	ids := make(map[string]SnakeId)
	lengths := make(map[SnakeId]uint8)
	healths := make(map[SnakeId]int8)
	heads := make(map[SnakeId]uint16)

	curSnakeId := 1
	ids[state.You.ID] = SnakeId(curSnakeId)
	for _, s := range state.Board.Snakes {
		if s.ID != state.You.ID {
			curSnakeId += 1
			ids[s.ID] = SnakeId(curSnakeId)
		}
	}

	list := []Tile{}
	for i := uint16(0); i < listLength; i++ {
		list = append(list, CreateEmptyTile())
	}

	for _, s := range state.Board.Snakes {
		length := len(s.Body)
		head := Point{X: int8(s.Body[0].X), Y: int8(s.Body[0].Y)}
		tail := Point{X: int8(s.Body[length-1].X), Y: int8(s.Body[length-1].Y)}
		headIndex := pointToIndex(head, width)
		headId := ids[s.ID]

		lengths[headId] = uint8(length)
		healths[headId] = int8(s.Health)
		heads[headId] = headIndex

		for i := length - 1; i > 0; i-- {
			bodyPart := Point{X: int8(s.Body[i].X), Y: int8(s.Body[i].Y)}
			nextBodyPart := Point{X: int8(s.Body[i-1].X), Y: int8(s.Body[i-1].Y)}
			fmt.Println("bodyPart ", bodyPart)
			fmt.Println("width ", width)
			index := pointToIndex(bodyPart, width)
			fmt.Println(index)
			tile := list[index]
			tile.SetBodyPart(headId, pointToIndex(nextBodyPart, width))
			list[index] = tile
		}

		headTile := list[headIndex]
		headTile.SetHead(headId, pointToIndex(tail, width))
		list[headIndex] = headTile
	}

	for _, h := range state.Board.Hazards {
		p := Point{X: int8(h.X), Y: int8(h.Y)}
		index := pointToIndex(p, width)
		tile := list[index]
		tile.SetHazard()
		list[index] = tile
	}

	for _, f := range state.Board.Food {
		p := Point{X: int8(f.X), Y: int8(f.Y)}
		index := pointToIndex(p, width)
		tile := list[index]
		tile.SetFood()
		list[index] = tile
	}

	return FastBoard{
		width:     width,
		height:    height,
		list:      list,
		ids:       ids,
		Lengths:   lengths,
		Healths:   healths,
		Heads:     heads,
		isWrapped: state.Game.Ruleset.Name == "wrapped",
	}
}

func (b *FastBoard) AdvanceBoard(moves []SnakeMove) {
	newHeads := make(map[SnakeId]uint16, len(moves))
	for _, m := range moves {
		headIndex := b.Heads[m.Id]
		newHeads[m.Id] = indexInDirection(m.Dir, headIndex, b.width, b.height, b.isWrapped)
	}

	deadSnakes := make(map[SnakeId]struct{}, len(newHeads))

	// do damage food and stuff
	for id, index := range newHeads {
		// head to head with larger or equal snake
		for oId, oIndex := range newHeads {
			if oIndex == index && oId != id && b.Lengths[oId] >= b.Lengths[id] {
				deadSnakes[id] = struct{}{}
			}
		}

		moveTile := b.list[index]

		// reduce health
		b.Healths[id] -= 1

		// feed snakes
		if moveTile.IsFood() {
			b.Healths[id] = 100
		}

		if moveTile.IsHazard() {
			// TODO: handle other types of hazard damage
			b.Healths[id] -= 100
		}

		// snake collision
		if moveTile.IsSnakeSegment() {
			snakeTailIndex := b.list[b.Heads[moveTile.id]].GetIdx()
			isNonTailSegment := index != snakeTailIndex
			didSnakeEat := b.list[newHeads[moveTile.id]].IsFood()

			if isNonTailSegment || didSnakeEat {
				deadSnakes[id] = struct{}{}
			}
		}

		if b.isOffBoard(indexToPoint(index, b.width)) {
			deadSnakes[id] = struct{}{}
		}

		// check for out of health
		if b.Healths[id] <= 0 {
			deadSnakes[id] = struct{}{}
		}
	}

	//  kill snakes
	fmt.Println("before kill")
	for id := range deadSnakes {
		b.kill(id)
	}
	fmt.Println("after kill")

	// move tails
	for id, index := range newHeads {
		if b.Healths[id] < 1 {
			continue
		}

		tailIndex := b.list[b.Heads[id]].GetIdx()
		didSnakeEat := b.isTileFood(index)

		if didSnakeEat {
			continue
		}

		//  update head with new tail
		b.setTileSnakeHead(b.Heads[id], id, b.list[tailIndex].GetIdx())
		// clear tail spot
		b.clearTile(tailIndex)
	}

	// move heads
	for id, index := range newHeads {
		if b.Healths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.list[oldHeadIndex].GetIdx()
		b.setTileSnakeHead(index, id, tailIndex)
		b.Heads[id] = index
		//  update neck with new head
		b.setTileSnakeBodyPart(oldHeadIndex, id, index)
	}
}

func (b *FastBoard) kill(id SnakeId) {
	headIndex := b.Heads[id]
	head := b.list[headIndex]
	tailIndex := head.GetIdx()

	currentIndex := tailIndex
	for currentIndex != headIndex {
		nextIndex := b.list[currentIndex].GetIdx()
		b.clearTile(currentIndex)
		currentIndex = nextIndex
	}

	b.clearTile(headIndex)
	b.Lengths[id] = 0
	b.Healths[id] = 0
	b.Heads[id] = 0
}

func (b *FastBoard) RandomRollout() {
	for !b.IsGameOver() {
		var moves []SnakeMove
		for id, l := range b.Lengths {
			fmt.Println("head", id, b.Heads[id])
			fmt.Println("length", id, b.Lengths[id])
			fmt.Println("health", id, b.Healths[id])
			if l == 0 {
				continue
			}

			sMoves := b.GetMovesForSnake(id)
			if len(sMoves) < 1 {
				sMoves = []SnakeMove{{Id: id, Dir: Left}}
			}
			randomMove := sMoves[rand.Intn(len(sMoves))]
			moves = append(moves, randomMove)
		}
		fmt.Println("advancing board")
		fmt.Println(b)
		b.AdvanceBoard(moves)
		fmt.Println("done advancing board")
	}
}

func (b *FastBoard) GetMovesForSnake(id SnakeId) []SnakeMove {
	moves := [4]Move{Up, Left, Right, Down}
	var possibleMoves []SnakeMove
	snakeHeadIndex := b.Heads[id]

	for _, dir := range moves {
		dirPoint := b.pointInDirection(dir, snakeHeadIndex)
		isOutOfBounds := b.isOffBoard(dirPoint)
		if isOutOfBounds {
			continue
		}

		dirIndex := pointToIndex(dirPoint, b.width)
		if b.isTileHazard(dirIndex) {
			continue
		}

		if b.isTileSnakeSegment(dirIndex) && !b.isTileSafeTail(dirIndex) {
			continue
		}

		possibleMoves = append(possibleMoves, SnakeMove{Id: id, Dir: dir})
	}

	return possibleMoves
}

func (b *FastBoard) setTileSnakeDoubleStack(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.list[index]
	tile.SetDoubleStack(id, nextIdx)
	b.list[index] = tile
}

func (b *FastBoard) setTileSnakeBodyPart(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.list[index]
	tile.SetBodyPart(id, nextIdx)
	b.list[index] = tile
}

func (b *FastBoard) setTileSnakeHead(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.list[index]
	tile.SetHead(id, nextIdx)
	b.list[index] = tile
}

func (b *FastBoard) clearTile(index uint16) {
	tile := b.list[index]
	tile.Clear()
	b.list[index] = tile
}

func (b *FastBoard) isTileHazard(index uint16) bool {
	return b.list[index].IsHazard()
}

func (b *FastBoard) isTileFood(index uint16) bool {
	return b.list[index].IsFood()
}

func (b *FastBoard) isTileSnakeHead(index uint16) bool {
	return b.list[index].IsSnakeHead()
}

func (b *FastBoard) isTileSnakeSegment(index uint16) bool {
	return b.list[index].IsSnakeSegment()
}

func (b *FastBoard) isTileSafeTail(index uint16) bool {
	id, ok := b.getSnakeIdAtTile(index)
	if ok {
		headIndex := b.Heads[id]
		tailIndex := b.list[headIndex].GetIdx()

		return tailIndex == index
	}

	return false
}

func (b *FastBoard) getSnakeIdAtTile(index uint16) (SnakeId, bool) {
	return b.list[index].GetSnakeId()
}

func (b *FastBoard) isOffBoard(p Point) bool {
	return p.X < 0 || p.X >= int8(b.width) || p.Y < 0 || p.Y >= int8(b.height)
}

func (b *FastBoard) IsGameOver() bool {
	var snakesLeft []SnakeId
	for id, health := range b.Healths {
		if health > 0 {
			snakesLeft = append(snakesLeft, id)
		}
	}

	return len(snakesLeft) < 2
}

func (b *FastBoard) pointInDirection(m Move, cur uint16) Point {
	p := addPoints(indexToPoint(cur, b.width), moveToPoint(m))
	if b.isWrapped {
		p = adjustForWrapped(p, b.width, b.height)
	}
	return p
}

//type FastBoard struct {
//list      []Tile
//ids       map[string]SnakeId
//Heads     map[SnakeId]uint16
//Lengths   map[SnakeId]uint8
//Healths   map[SnakeId]int8
//width     uint16
//height    uint16
//isWrapped bool
//}

func (b *FastBoard) Clone() FastBoard {
	newBoard := FastBoard{}

	newBoard.list = make([]Tile, len(b.list))
	copy(newBoard.list, b.list)

	newBoard.ids = make(map[string]SnakeId)
	for sId, id := range b.ids {
		if b.Lengths[id] > 0 {
			newBoard.ids[sId] = id
		}
	}

	newBoard.Heads = make(map[SnakeId]uint16)
	for id, h := range b.Heads {
		if b.Lengths[id] > 0 {
			newBoard.Heads[id] = h
		}
	}

	newBoard.Lengths = make(map[SnakeId]uint8)
	for id, l := range b.Lengths {
		if b.Lengths[id] > 0 {
			newBoard.Lengths[id] = l
		}
	}

	newBoard.Healths = make(map[SnakeId]int8)
	for id, h := range b.Healths {
		if b.Lengths[id] > 0 {
			newBoard.Healths[id] = h
		}
	}

	newBoard.width = b.width
	newBoard.height = b.height
	newBoard.isWrapped = b.isWrapped

	return newBoard
}
