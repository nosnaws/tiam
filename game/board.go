package game

import "fmt"

type board struct {
	list    []Tile
	ids     map[string]SnakeId
	heads   map[SnakeId]uint16
	lengths map[SnakeId]uint8
	healths map[SnakeId]int8
	width   uint16
	height  uint16
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
		width:   width,
		height:  height,
		list:    list,
		ids:     ids,
		lengths: lengths,
		healths: healths,
		heads:   heads,
	}
}

func (b *board) advanceBoard(moves []snakeMove) {
	newHeads := make(map[SnakeId]uint16, len(moves))
	for _, m := range moves {
		headIndex := b.heads[m.id]
		newHeads[m.id] = indexInDirection(m.dir, headIndex, b.width)
	}

	deadSnakes := make(map[SnakeId]struct{}, len(newHeads))
	// head to head with larger or equal snake
	for id, index := range newHeads {
		for oId, oIndex := range newHeads {
			if oIndex == index && oId != id && b.lengths[oId] >= b.lengths[id] {
				deadSnakes[id] = struct{}{}
			}
		}
	}

	for id, index := range newHeads {
		if _, ok := deadSnakes[id]; ok {
			continue
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
			didSnakeMove := newHeads[moveTile.id] == b.heads[moveTile.id]

			fmt.Println("Found snake body", snakeTailIndex)
			if didSnakeMove == true || index != snakeTailIndex || b.list[newHeads[moveTile.id]].IsFood() {
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

		if _, ok := deadSnakes[id]; ok {
			continue
		}

		// move snake!
		var tailIndex uint16
		if moveTile.IsFood() {
			tailIndex = b.list[b.heads[id]].GetIdx()
		} else {
			oldTailIndex := b.list[b.heads[id]].GetIdx()
			tailIndex = b.list[oldTailIndex].GetIdx()

			// remove tail tile since snake did not eat
			b.list[oldTailIndex].Clear()
		}

		// update neck to point to new head
		oldHeadIndex := b.heads[id]
		b.setTileSnakeBodyPart(oldHeadIndex, id, index)

		// move snake head
		b.setTileSnakeHead(index, id, tailIndex)
		b.heads[id] = index
	}

	for id := range deadSnakes {
		b.kill(id)
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
