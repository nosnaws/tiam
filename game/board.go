package game

import "fmt"

type board struct {
	list      []Tile
	ids       map[string]SnakeId
	heads     map[SnakeId]uint16
	lengths   map[SnakeId]uint8
	healths   map[SnakeId]int8
	width     uint16
	height    uint16
	isWrapped bool
}

type snakeMove struct {
	id  SnakeId
	dir Move
}

type idToNewHead struct {
	id    SnakeId
	index uint16
}

func BuildBoard(state GameState) board {
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
			index := pointToIndex(bodyPart, width)
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

	return board{
		width:     width,
		height:    height,
		list:      list,
		ids:       ids,
		lengths:   lengths,
		healths:   healths,
		heads:     heads,
		isWrapped: state.Game.Ruleset.Name == "wrapped",
	}
}

func (b *board) advanceBoard(moves []snakeMove) {
	newHeads := make(map[SnakeId]uint16, len(moves))
	for _, m := range moves {
		headIndex := b.heads[m.id]
		newHeads[m.id] = indexInDirection(m.dir, headIndex, b.width, b.height, b.isWrapped)
	}

	deadSnakes := make(map[SnakeId]struct{}, len(newHeads))

	// do damage food and stuff
	for id, index := range newHeads {
		// head to head with larger or equal snake
		for oId, oIndex := range newHeads {
			if oIndex == index && oId != id && b.lengths[oId] >= b.lengths[id] {
				deadSnakes[id] = struct{}{}
			}
		}

		moveTile := b.list[index]

		// reduce health
		b.healths[id] -= 1

		// feed snakes
		if moveTile.IsFood() {
			b.healths[id] = 100
		}

		if moveTile.IsHazard() {
			// TODO: handle other types of hazard damage
			b.healths[id] -= 100
		}

		// snake collision
		if moveTile.IsSnakeSegment() {
			snakeTailIndex := b.list[b.heads[moveTile.id]].GetIdx()
			isNonTailSegment := index != snakeTailIndex
			didSnakeEat := b.list[newHeads[moveTile.id]].IsFood()

			fmt.Println("Found snake body", snakeTailIndex)
			fmt.Println("will snake eat ", b.list[newHeads[moveTile.id]].IsFood())
			if isNonTailSegment || didSnakeEat {
				fmt.Println("Snake collision", id)
				deadSnakes[id] = struct{}{}
			}
		}

		if b.isOffBoard(indexToPoint(index, b.width)) {
			fmt.Println("Off board", id)
			deadSnakes[id] = struct{}{}
		}

		// check for out of health
		if b.healths[id] <= 0 {
			fmt.Println(b.healths)
			fmt.Println("Out of health", id)
			deadSnakes[id] = struct{}{}
		}
	}

	//  kill snakes
	for id := range deadSnakes {
		b.kill(id)
	}

	// move tails
	for id, index := range newHeads {
		if b.healths[id] < 1 {
			continue
		}

		tailIndex := b.list[b.heads[id]].GetIdx()
		didSnakeEat := b.isTileFood(index)

		if didSnakeEat {
			continue
		}

		//  update head with new tail
		b.setTileSnakeHead(b.heads[id], id, b.list[tailIndex].GetIdx())
		// clear tail spot
		b.clearTile(tailIndex)
	}

	// move heads
	for id, index := range newHeads {
		if b.healths[id] < 1 {
			continue
		}

		oldHeadIndex := b.heads[id]
		tailIndex := b.list[oldHeadIndex].GetIdx()
		b.setTileSnakeHead(index, id, tailIndex)
		b.heads[id] = index
		//  update neck with new head
		b.setTileSnakeBodyPart(oldHeadIndex, id, index)
	}
}

func (b *board) kill(id SnakeId) {
	headIndex := b.heads[id]
	head := b.list[headIndex]
	tailIndex := head.GetIdx()

	currentIndex := tailIndex
	for currentIndex != headIndex {
		nextIndex := b.list[currentIndex].GetIdx()
		b.clearTile(currentIndex)
		currentIndex = nextIndex
	}

	b.clearTile(headIndex)
	b.lengths[id] = 0
	b.healths[id] = 0
	b.heads[id] = 0
}

func (b *board) getMovesForSnake(id SnakeId) []snakeMove {
	moves := []Move{Up, Left, Right, Down}
	var possibleMoves []snakeMove
	snakeHead := b.heads[id]

	for _, dir := range moves {
		dirIndex := indexInDirection(dir, snakeHead, b.width, b.height, b.isWrapped)
		isHazard := b.list[dirIndex].IsHazard()
		if isHazard {
			continue
		}

		isOutOfBounds := b.isOffBoard(indexToPoint(dirIndex, b.width))
		if isOutOfBounds {
			continue
		}

		possibleMoves = append(possibleMoves, snakeMove{id: id, dir: dir})
	}

	return possibleMoves
}

func (b *board) setTileSnakeDoubleStack(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.list[index]
	tile.SetDoubleStack(id, nextIdx)
	b.list[index] = tile
}

func (b *board) setTileSnakeBodyPart(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.list[index]
	tile.SetBodyPart(id, nextIdx)
	b.list[index] = tile
}

func (b *board) setTileSnakeHead(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.list[index]
	tile.SetHead(id, nextIdx)
	b.list[index] = tile
}

func (b *board) clearTile(index uint16) {
	tile := b.list[index]
	tile.Clear()
	b.list[index] = tile
}

func (b *board) isTileHazard(index uint16) bool {
	return b.list[index].IsHazard()
}

func (b *board) isTileFood(index uint16) bool {
	return b.list[index].IsFood()
}

func (b *board) isTileSnakeHead(index uint16) bool {
	return b.list[index].IsSnakeHead()
}

func (b *board) isTileSnakeBody(index uint16) bool {
	return b.list[index].IsSnakeBody()
}

func (b *board) getSnakeIdAtTile(index uint16) (SnakeId, bool) {
	return b.list[index].GetSnakeId()
}

func (b *board) isOffBoard(p Point) bool {
	return p.X < 0 || p.X >= int8(b.width) || p.Y < 0 || p.Y >= int8(b.height)
}

func (b *board) isGameOver() bool {
	var snakesLeft []SnakeId
	for id, health := range b.healths {
		if health > 0 {
			snakesLeft = append(snakesLeft, id)
		}
	}

	return len(snakesLeft) < 2
}
