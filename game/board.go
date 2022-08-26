package game

import (
	"fmt"
	"math/rand"
)

const MeId = SnakeId(1)

var allMoves = [4]Move{Up, Left, Right, Down}

type idsMap map[string]SnakeId
type headsMap map[SnakeId]uint16
type lengthsMap map[SnakeId]uint8
type healthsMap map[SnakeId]int8
type hazardDepthMap map[uint16]int8

type FastBoard struct {
	list         []Tile
	ids          idsMap
	Heads        headsMap
	Lengths      lengthsMap
	Healths      healthsMap
	width        uint16
	height       uint16
	isWrapped    bool
	hazardDamage int8
	hazardDepth  hazardDepthMap
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
	hazardDamage := int8(state.Game.Ruleset.Settings.HazardDamagePerTurn)
	listLength := height * width

	ids := make(idsMap)
	lengths := make(lengthsMap)
	healths := make(healthsMap)
	heads := make(headsMap)
	hazardDepth := make(hazardDepthMap)

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

		if isSnakeInitialState(s) {
			list[headIndex].SetTripleStack(headId)
		} else if isSnakeDoubleButt(s) {

			doubleStackP := Point{X: int8(s.Body[length-1].X), Y: int8(s.Body[length-1].Y)}
			doubleStackIndex := pointToIndex(doubleStackP, width)
			nextBodyPart := Point{X: int8(s.Body[length-3].X), Y: int8(s.Body[length-3].Y)}
			list[doubleStackIndex].SetDoubleStack(headId, pointToIndex(nextBodyPart, width))

			for i := length - 3; i > 0; i-- {
				bodyPart := Point{X: int8(s.Body[i].X), Y: int8(s.Body[i].Y)}
				nextBodyPart := Point{X: int8(s.Body[i-1].X), Y: int8(s.Body[i-1].Y)}
				index := pointToIndex(bodyPart, width)
				tile := list[index]
				tile.SetBodyPart(headId, pointToIndex(nextBodyPart, width))
				list[index] = tile
			}

			list[headIndex].SetHead(headId, pointToIndex(tail, width))
		} else {

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
	}

	for _, h := range state.Board.Hazards {
		p := Point{X: int8(h.X), Y: int8(h.Y)}
		index := pointToIndex(p, width)
		tile := list[index]
		tile.SetHazard()
		hazardDepth[index] += 1
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
		width:        width,
		height:       height,
		list:         list,
		ids:          ids,
		Lengths:      lengths,
		Healths:      healths,
		Heads:        heads,
		isWrapped:    state.Game.Ruleset.Name == "wrapped",
		hazardDamage: hazardDamage,
		hazardDepth:  hazardDepth,
	}
}

func (b *FastBoard) AdvanceBoard(moves map[SnakeId]Move) {
	deadSnakes := make([]SnakeId, len(b.Heads))

	// do damage food and stuff
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.width, b.height, b.isWrapped)
		newIndex := pointToIndex(newPoint, b.width)
		isDead := false

		// head to head with larger or equal snake
		for oId, oDir := range moves {
			oPoint := pointInDirection(oDir, b.Heads[oId], b.width, b.height, b.isWrapped)
			oIndex := pointToIndex(oPoint, b.width)

			if oIndex == newIndex && oId != id && b.Lengths[oId] >= b.Lengths[id] {
				isDead = true
			}
		}

		if b.isOffBoard(newPoint) {
			isDead = true
		}

		if !isDead {
			moveTile := b.list[newIndex]

			// reduce health
			b.Healths[id] -= 1

			// feed snakes
			if moveTile.IsFood() {
				b.Healths[id] = 100
			}

			if moveTile.IsHazard() {
				b.Healths[id] -= b.hazardDamage * b.hazardDepth[newIndex]
			}

			// snake collision
			if moveTile.IsSnakeSegment() {
				snakeTailIndex := b.list[b.Heads[moveTile.id]].GetIdx()
				isNonTailSegment := newIndex != snakeTailIndex
				didSnakeEat := false

				snakeDir := moves[moveTile.id]
				snakeNewHead := pointInDirection(snakeDir, b.Heads[moveTile.id], b.width, b.height, b.isWrapped)
				if !b.isOffBoard(snakeNewHead) {
					didSnakeEat = b.list[pointToIndex(snakeNewHead, b.width)].IsFood()
				}

				if isNonTailSegment || didSnakeEat {
					isDead = true
				}
			}

			// check for out of health
			if b.Healths[id] <= 0 {
				isDead = true
			}
		}

		if isDead {
			deadSnakes = append(deadSnakes, id)
		}
	}

	//  kill snakes
	for _, id := range deadSnakes {
		if id != 0 {
			b.kill(id)
		}
	}

	// move tails
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.width, b.height, b.isWrapped)
		newIndex := pointToIndex(newPoint, b.width)

		if b.Lengths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.list[oldHeadIndex].GetIdx()
		didSnakeEat := b.IsTileFood(newIndex)

		oldHeadTile := b.list[oldHeadIndex]
		tailTile := b.list[tailIndex]
		if oldHeadTile.IsTripleStack() || tailTile.IsDoubleStack() {
			continue
		}

		if didSnakeEat {
			b.Lengths[id] += 1
			continue
		}

		//  update head with new tail
		b.setTileSnakeHead(b.Heads[id], id, b.list[tailIndex].GetIdx())
		// clear tail spot
		b.clearTile(tailIndex)
	}

	// move heads
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.width, b.height, b.isWrapped)
		newIndex := pointToIndex(newPoint, b.width)

		if b.Healths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.list[oldHeadIndex].GetIdx()

		if b.list[oldHeadIndex].IsTripleStack() {
			b.setTileSnakeHead(newIndex, id, oldHeadIndex)
			b.setTileSnakeDoubleStack(oldHeadIndex, id, newIndex)
		} else if b.list[tailIndex].IsDoubleStack() {
			b.setTileSnakeHead(newIndex, id, tailIndex)
			b.setTileSnakeBodyPart(oldHeadIndex, id, newIndex)
			b.setTileSnakeBodyPart(tailIndex, id, oldHeadIndex)
		} else {
			b.setTileSnakeHead(newIndex, id, tailIndex)
			//  update neck with new head
			b.setTileSnakeBodyPart(oldHeadIndex, id, newIndex)
		}
		b.Heads[id] = newIndex
	}
}

func (b *FastBoard) AdvanceBoardTurnBased(moves map[SnakeId]Move) {
	deadSnakes := make([]SnakeId, len(b.Heads))

	// do damage food and stuff
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.width, b.height, b.isWrapped)
		newIndex := pointToIndex(newPoint, b.width)
		isDead := false

		// head to head with larger or equal snake
		if b.isTileSnakeHead(newIndex) {
			sId, ok := b.getSnakeIdAtTile(newIndex)
			if ok && b.Lengths[sId] >= b.Lengths[id] {
				isDead = true
			}
		}
		//for oId, oDir := range moves {
		//oPoint := pointInDirection(oDir, b.Heads[oId], b.width, b.height, b.isWrapped)
		//oIndex := pointToIndex(oPoint, b.width)

		//if oIndex == newIndex && oId != id && b.Lengths[oId] >= b.Lengths[id] {
		//isDead = true
		//}
		//}

		if b.isOffBoard(newPoint) {
			isDead = true
		}

		if !isDead {
			moveTile := b.list[newIndex]

			// reduce health
			b.Healths[id] -= 1

			// feed snakes
			if moveTile.IsFood() {
				b.Healths[id] = 100
			}

			if moveTile.IsHazard() {
				b.Healths[id] -= b.hazardDamage * b.hazardDepth[newIndex]
			}

			// snake collision
			if moveTile.IsSnakeSegment() {
				snakeTailIndex := b.list[b.Heads[moveTile.id]].GetIdx()
				isNonTailSegment := newIndex != snakeTailIndex
				didSnakeEat := false

				snakeDir := moves[moveTile.id]
				snakeNewHead := pointInDirection(snakeDir, b.Heads[moveTile.id], b.width, b.height, b.isWrapped)
				if !b.isOffBoard(snakeNewHead) {
					didSnakeEat = b.list[pointToIndex(snakeNewHead, b.width)].IsFood()
				}

				if isNonTailSegment || didSnakeEat {
					isDead = true
				}
			}

			// check for out of health
			if b.Healths[id] <= 0 {
				isDead = true
			}
		}

		if isDead {
			deadSnakes = append(deadSnakes, id)
		}
	}

	//  kill snakes
	for _, id := range deadSnakes {
		if id != 0 {
			b.kill(id)
		}
	}

	// move tails
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.width, b.height, b.isWrapped)
		newIndex := pointToIndex(newPoint, b.width)

		if b.Lengths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.list[oldHeadIndex].GetIdx()
		didSnakeEat := b.IsTileFood(newIndex)

		oldHeadTile := b.list[oldHeadIndex]
		tailTile := b.list[tailIndex]
		if oldHeadTile.IsTripleStack() || tailTile.IsDoubleStack() {
			continue
		}

		if didSnakeEat {
			b.Lengths[id] += 1
			continue
		}

		//  update head with new tail
		b.setTileSnakeHead(b.Heads[id], id, b.list[tailIndex].GetIdx())
		// clear tail spot
		b.clearTile(tailIndex)
	}

	// move heads
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.width, b.height, b.isWrapped)
		newIndex := pointToIndex(newPoint, b.width)

		if b.Healths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.list[oldHeadIndex].GetIdx()

		if b.list[oldHeadIndex].IsTripleStack() {
			b.setTileSnakeHead(newIndex, id, oldHeadIndex)
			b.setTileSnakeDoubleStack(oldHeadIndex, id, newIndex)
		} else if b.list[tailIndex].IsDoubleStack() {
			b.setTileSnakeHead(newIndex, id, tailIndex)
			b.setTileSnakeBodyPart(oldHeadIndex, id, newIndex)
			b.setTileSnakeBodyPart(tailIndex, id, oldHeadIndex)
		} else {
			b.setTileSnakeHead(newIndex, id, tailIndex)
			//  update neck with new head
			b.setTileSnakeBodyPart(oldHeadIndex, id, newIndex)
		}
		b.Heads[id] = newIndex
	}
}

func (b *FastBoard) kill(id SnakeId) {
	headIndex := b.Heads[id]
	head := b.list[headIndex]
	tailIndex := head.GetIdx()

	if b.list[headIndex].IsTripleStack() {
		b.clearTile(headIndex)
	} else if b.list[tailIndex].IsDoubleStack() {
		b.clearTile(tailIndex)
		b.clearTile(headIndex)
	} else {
		currentIndex := tailIndex
		for currentIndex != headIndex {
			nextIndex := b.list[currentIndex].GetIdx()
			b.clearTile(currentIndex)
			currentIndex = nextIndex
		}
		b.clearTile(headIndex)
	}

	b.Lengths[id] = 0
	b.Healths[id] = 0
	b.Heads[id] = 0
}

func (b *FastBoard) RandomRollout() {
	//turn := 0
	moves := make(map[SnakeId]Move, len(b.Lengths))
	for !b.IsGameOver() {
		//moves := make([]SnakeMove, 0, len(b.Lengths))

		for id, l := range b.Lengths {
			if l == 0 {
				moves[id] = ""
				continue
			}

			sMoves := b.GetMovesForSnake(id)
			randomMove := sMoves[rand.Intn(len(sMoves))]
			//moves = append(moves, randomMove)
			moves[id] = randomMove.Dir
		}
		//turn += 1
		b.AdvanceBoard(moves)
	}
}

func (b *FastBoard) GetMovesForSnake(id SnakeId) []SnakeMove {
	var possibleMoves []SnakeMove
	snakeHeadIndex := b.Heads[id]

	moves := b.GetNeighbors(snakeHeadIndex)
	for _, m := range moves {
		possibleMoves = append(possibleMoves, SnakeMove{Id: id, Dir: m})
	}
	// No moves, go left
	if len(possibleMoves) == 0 {
		possibleMoves = append(possibleMoves, SnakeMove{Id: id, Dir: Left})
	}
	return possibleMoves
}

func (b *FastBoard) GetNeighbors(index uint16) []Move {
	possibleMoves := make([]Move, 0, 4)

	for _, dir := range allMoves {
		dirPoint := b.pointInDirection(dir, index)
		isOutOfBounds := b.isOffBoard(dirPoint)
		if isOutOfBounds {
			continue
		}

		dirIndex := pointToIndex(dirPoint, b.width)
		if b.isTileHazard(dirIndex) && b.hazardDamage >= 100 {
			continue
		}

		if b.isTileSnakeSegment(dirIndex) && !b.isTileSafeTail(dirIndex) {
			continue
		}

		possibleMoves = append(possibleMoves, dir)
	}

	return possibleMoves
}

func (b *FastBoard) setTileSnakeTripleStack(index uint16, id SnakeId) {
	tile := b.list[index]
	tile.SetTripleStack(id)
	b.list[index] = tile
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

func (b *FastBoard) IsTileFood(index uint16) bool {
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

		return tailIndex == index && !b.list[tailIndex].IsDoubleStack()
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
	snakesLeft := 0
	for id := range b.Lengths {
		if b.IsSnakeAlive(id) {
			snakesLeft += 1
		}
	}

	return snakesLeft < 2
}

func (b *FastBoard) IsSnakeAlive(id SnakeId) bool {
	if b.Lengths[id] > 0 {
		return true
	}

	return false
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

	newBoard.ids = make(idsMap, len(b.ids))
	for sId, id := range b.ids {
		newBoard.ids[sId] = id
	}

	newBoard.Heads = make(headsMap, len(b.Heads))
	for id, h := range b.Heads {
		newBoard.Heads[id] = h
	}

	newBoard.Lengths = make(lengthsMap, len(b.Lengths))
	for id, l := range b.Lengths {
		newBoard.Lengths[id] = l
	}

	newBoard.Healths = make(healthsMap, len(b.Healths))
	for id, h := range b.Healths {
		newBoard.Healths[id] = h
	}

	newBoard.hazardDepth = make(hazardDepthMap, len(b.hazardDepth))
	for index, depth := range b.hazardDepth {
		newBoard.hazardDepth[index] = depth
	}

	newBoard.width = b.width
	newBoard.height = b.height
	newBoard.isWrapped = b.isWrapped
	newBoard.hazardDamage = b.hazardDamage

	return newBoard
}

func (b *FastBoard) tileToString(index uint16) string {
	if b.isTileSnakeHead(index) {
		return fmt.Sprintf(" %dh ", b.list[index].id)
	}
	if b.isTileSnakeSegment(index) {
		return fmt.Sprintf(" %ds ", b.list[index].id)
	}
	if b.IsTileFood(index) && b.isTileHazard(index) {
		return " fz "
	}
	if b.IsTileFood(index) {
		return " ff "
	}
	if b.isTileHazard(index) {
		return " zz "
	}
	return " __ "
}

func (b *FastBoard) Print() {
	fmt.Println("######")
	for id, length := range b.Lengths {
		fmt.Printf("%d - length:%d\n", id, length)
	}
	for id, health := range b.Healths {
		fmt.Printf("%d - health:%d\n", id, health)
	}
	for y := int(b.height - 1); y >= 0; y-- {
		var line string
		for x := 0; x < int(b.width); x++ {
			p := Point{X: int8(x), Y: int8(y)}
			line = line + b.tileToString(pointToIndex(p, b.width))
		}
		fmt.Println(line)
	}
}
