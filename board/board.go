package board

import (
	"fmt"

	b "github.com/nosnaws/tiam/battlesnake"
)

const MeId = SnakeId(1)

var allMoves = [4]Move{Up, Left, Right, Down}

type idsMap map[string]SnakeId
type headsMap map[SnakeId]uint16
type lengthsMap map[SnakeId]uint8
type healthsMap map[SnakeId]int8
type hazardDepthMap map[uint16]int8

type FastBoard struct {
	List         []Tile
	Ids          idsMap
	Heads        headsMap
	Lengths      lengthsMap
	Healths      healthsMap
	Width        uint16
	Height       uint16
	IsWrapped    bool
	hazardDamage int8
	hazardDepth  hazardDepthMap
	hasEaten     map[SnakeId]bool
}

type SnakeMove struct {
	Id  SnakeId
	Dir Move
}

type idToNewHead struct {
	id    SnakeId
	index uint16
}

func BuildBoard(state b.GameState) FastBoard {
	height := uint16(state.Board.Height)
	width := uint16(state.Board.Width)
	hazardDamage := int8(state.Game.Ruleset.Settings.HazardDamagePerTurn)
	listLength := height * width

	ids := make(idsMap)
	lengths := make(lengthsMap)
	healths := make(healthsMap)
	heads := make(headsMap)
	hazardDepth := make(hazardDepthMap)
	hasEaten := make(map[SnakeId]bool)

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
		Width:        width,
		Height:       height,
		List:         list,
		Ids:          ids,
		Lengths:      lengths,
		Healths:      healths,
		Heads:        heads,
		IsWrapped:    state.Game.Ruleset.Name == "wrapped",
		hazardDamage: hazardDamage,
		hazardDepth:  hazardDepth,
		hasEaten:     hasEaten,
	}
}

func (b *FastBoard) AdvanceBoard(moves map[SnakeId]Move) {
	deadSnakes := make([]SnakeId, len(b.Heads))

	// do damage food and stuff
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.Width, b.Height, b.IsWrapped)
		newIndex := pointToIndex(newPoint, b.Width)
		isDead := false

		// head to head with larger or equal snake
		for oId, oDir := range moves {
			oPoint := pointInDirection(oDir, b.Heads[oId], b.Width, b.Height, b.IsWrapped)
			oIndex := pointToIndex(oPoint, b.Width)

			if oIndex == newIndex && oId != id && b.Lengths[oId] >= b.Lengths[id] {
				isDead = true
			}
		}

		if b.isOffBoard(newPoint) {
			isDead = true
		}

		if !isDead {
			moveTile := b.List[newIndex]

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
				snakeTailIndex := b.List[b.Heads[moveTile.id]].GetIdx()
				isNonTailSegment := newIndex != snakeTailIndex
				didSnakeEat := false

				snakeDir := moves[moveTile.id]
				snakeNewHead := pointInDirection(snakeDir, b.Heads[moveTile.id], b.Width, b.Height, b.IsWrapped)
				if !b.isOffBoard(snakeNewHead) {
					didSnakeEat = b.List[pointToIndex(snakeNewHead, b.Width)].IsFood()
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

	eatingSnakes := make(map[SnakeId]bool)

	// move tails
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.Width, b.Height, b.IsWrapped)
		newIndex := pointToIndex(newPoint, b.Width)

		if b.Lengths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.List[oldHeadIndex].GetIdx()
		didSnakeEat := b.IsTileFood(newIndex)

		oldHeadTile := b.List[oldHeadIndex]
		tailTile := b.List[tailIndex]
		if oldHeadTile.IsTripleStack() {
			continue
		}

		next := b.List[tailIndex].GetIdx()
		if didSnakeEat {
			eatingSnakes[id] = true
			b.Lengths[id] += 1

			if !tailTile.IsDoubleStack() {
				b.setTileSnakeDoubleStack(next, id, b.List[next].GetIdx())
				// clear tail spot
				b.clearTile(tailIndex)
			}
			//  update head with new tail
			b.setTileSnakeHead(b.Heads[id], id, next)
		} else {
			if !b.List[tailIndex].IsDoubleStack() {
				//  update head with new tail
				b.setTileSnakeHead(b.Heads[id], id, next)
				// clear tail spot
				b.clearTile(tailIndex)
			}
		}
	}

	// move heads
	for id, dir := range moves {
		if dir == "" {
			continue
		}

		newPoint := pointInDirection(dir, b.Heads[id], b.Width, b.Height, b.IsWrapped)
		newIndex := pointToIndex(newPoint, b.Width)

		if b.Healths[id] < 1 {
			continue
		}

		oldHeadIndex := b.Heads[id]
		tailIndex := b.List[oldHeadIndex].GetIdx()

		if b.List[oldHeadIndex].IsTripleStack() {
			b.setTileSnakeHead(newIndex, id, oldHeadIndex)
			b.setTileSnakeDoubleStack(oldHeadIndex, id, newIndex)
		} else if b.List[tailIndex].IsDoubleStack() && b.Lengths[id] > 3 && !eatingSnakes[id] {
			b.setTileSnakeHead(newIndex, id, tailIndex)
			b.setTileSnakeBodyPart(oldHeadIndex, id, newIndex)
			b.setTileSnakeBodyPart(tailIndex, id, b.List[tailIndex].GetIdx())
		} else if b.List[tailIndex].IsDoubleStack() && b.Lengths[id] == 3 {
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
	head := b.List[headIndex]
	tailIndex := head.GetIdx()

	if b.List[headIndex].IsTripleStack() {
		b.clearTile(headIndex)
	} else if b.List[tailIndex].IsDoubleStack() {
		b.clearTile(tailIndex)
		b.clearTile(headIndex)
	} else {
		currentIndex := tailIndex
		for currentIndex != headIndex {
			nextIndex := b.List[currentIndex].GetIdx()
			b.clearTile(currentIndex)
			currentIndex = nextIndex
		}
		b.clearTile(headIndex)
	}

	b.Lengths[id] = 0
	b.Healths[id] = 0
	b.Heads[id] = 0
}

func (b *FastBoard) GetMovesForSnake(id SnakeId) []SnakeMove {
	possibleMoves := b.GetMovesForSnakeNoDefault(id)
	// No moves, go left
	if len(possibleMoves) == 0 {
		possibleMoves = append(possibleMoves, SnakeMove{Id: id, Dir: Left})
	}
	return possibleMoves
}

func (b *FastBoard) GetAllFood() []uint16 {
	food := []uint16{}
	for i, t := range b.List {
		if t.IsFood() {
			food = append(food, uint16(i))
		}
	}
	return food
}

func (b *FastBoard) GetMovesForSnakeNoDefault(id SnakeId) []SnakeMove {
	var possibleMoves []SnakeMove
	snakeHeadIndex := b.Heads[id]

	moves := b.GetNeighbors(snakeHeadIndex)
	for _, m := range moves {
		possibleMoves = append(possibleMoves, SnakeMove{Id: id, Dir: m})
	}
	return possibleMoves
}

func (b *FastBoard) GetNeighborsUnsafe(index uint16) []Move {
	possibleMoves := make([]Move, 0, 4)

	for _, dir := range allMoves {
		dirPoint := b.pointInDirection(dir, index)
		isOutOfBounds := b.isOffBoard(dirPoint)
		if isOutOfBounds {
			continue
		}

		dirIndex := pointToIndex(dirPoint, b.Width)
		if b.isTileHazard(dirIndex) && b.hazardDamage >= 100 {
			continue
		}

		if b.isTileSnakeSegment(dirIndex) && !b.isTileUnSafeTail(dirIndex) {
			continue
		}

		possibleMoves = append(possibleMoves, dir)
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

		dirIndex := pointToIndex(dirPoint, b.Width)
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

func (b *FastBoard) GetNeighborIndices(index uint16) []uint16 {
	possibleMoves := make([]uint16, 0, 4)

	for _, dir := range allMoves {
		dirPoint := b.pointInDirection(dir, index)
		isOutOfBounds := b.isOffBoard(dirPoint)
		if isOutOfBounds {
			continue
		}

		dirIndex := pointToIndex(dirPoint, b.Width)
		if b.isTileHazard(dirIndex) && b.hazardDamage >= 100 {
			continue
		}

		if b.isTileSnakeSegment(dirIndex) && !b.isTileSafeTail(dirIndex) {
			continue
		}

		possibleMoves = append(possibleMoves, dirIndex)
	}

	return possibleMoves
}

func (b *FastBoard) GetHazardNeighbors(index uint16) []uint16 {
	neighbors := b.GetNeighbors(index)

	hNeighbors := []uint16{}
	for _, neighbor := range neighbors {
		nIndex := IndexInDirection(neighbor, index, b.Width, b.Height, b.IsWrapped)
		if b.isTileHazard(nIndex) {
			hNeighbors = append(hNeighbors, nIndex)
		}
	}
	return hNeighbors
}

func (b *FastBoard) GetSnakeTail(id SnakeId) uint16 {
	return b.List[b.Heads[id]].GetIdx()
}

func (b *FastBoard) GetLastMove(id SnakeId) Move {
	neckIndex := b.GetSnakeNeck(id)
	headIndex := b.Heads[id]

	if neckIndex == headIndex {
		return ""
	}

	if b.IndexInDirection(Left, neckIndex) == headIndex {
		return Left
	}
	if b.IndexInDirection(Right, neckIndex) == headIndex {
		return Right
	}
	if b.IndexInDirection(Up, neckIndex) == headIndex {
		return Up
	}
	return Down
}

func (b *FastBoard) IndexInDirection(move Move, current uint16) uint16 {
	return IndexInDirection(move, current, b.Width, b.Height, b.IsWrapped)
}

func (b *FastBoard) GetSnakeNeck(id SnakeId) uint16 {
	headIndex := b.Heads[id]
	tailIndex := b.List[headIndex].GetIdx()

	if b.List[headIndex].IsTripleStack() {
		return headIndex
	} else if b.List[tailIndex].IsDoubleStack() {
		return tailIndex
	} else {
		currentIndex := tailIndex
		lastIndex := currentIndex
		for currentIndex != headIndex {
			nextIndex := b.List[currentIndex].GetIdx()
			lastIndex = currentIndex
			currentIndex = nextIndex
		}
		return lastIndex
	}
}

func (b *FastBoard) setTileSnakeTripleStack(index uint16, id SnakeId) {
	tile := b.List[index]
	tile.SetTripleStack(id)
	b.List[index] = tile
}

func (b *FastBoard) setTileSnakeDoubleStack(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.List[index]
	tile.SetDoubleStack(id, nextIdx)
	b.List[index] = tile
}

func (b *FastBoard) setTileSnakeBodyPart(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.List[index]
	tile.SetBodyPart(id, nextIdx)
	b.List[index] = tile
}

func (b *FastBoard) setTileSnakeHead(index uint16, id SnakeId, nextIdx uint16) {
	tile := b.List[index]
	tile.SetHead(id, nextIdx)
	b.List[index] = tile
}

func (b *FastBoard) clearTile(index uint16) {
	tile := b.List[index]
	tile.Clear()
	b.List[index] = tile
}

func (b *FastBoard) isTileHazard(index uint16) bool {
	return b.List[index].IsHazard()
}

func (b *FastBoard) IsTileFood(index uint16) bool {
	return b.List[index].IsFood()
}

func (b *FastBoard) IsTileSnakeHead(index uint16) bool {
	return b.List[index].IsSnakeHead()
}

func (b *FastBoard) isTileSnakeSegment(index uint16) bool {
	return b.List[index].IsSnakeSegment()
}

func (b *FastBoard) isTileNonHeadSnakeSegment(index uint16) bool {
	return b.List[index].IsNonHeadSegment()
}

func (b *FastBoard) isTileSafeTail(index uint16) bool {
	id, ok := b.GetSnakeIdAtTile(index)
	if ok {
		headIndex := b.Heads[id]
		tailIndex := b.List[headIndex].GetIdx()

		return tailIndex == index && !b.List[tailIndex].IsDoubleStack()
	}

	return false
}

func (b *FastBoard) isTileUnSafeTail(index uint16) bool {
	id, ok := b.GetSnakeIdAtTile(index)
	if ok {
		headIndex := b.Heads[id]
		tailIndex := b.List[headIndex].GetIdx()

		return tailIndex == index
	}

	return false
}

func (b *FastBoard) GetSnakeIdAtTile(index uint16) (SnakeId, bool) {
	return b.List[index].GetSnakeId()
}

func (b *FastBoard) isOffBoard(p Point) bool {
	return p.X < 0 || p.X >= int8(b.Width) || p.Y < 0 || p.Y >= int8(b.Height)
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

func (b *FastBoard) RemoveSnake(id SnakeId) {
	b.kill(id)
}

func (b *FastBoard) IsSnakeAlive(id SnakeId) bool {
	if b.Lengths[id] > 0 {
		return true
	}

	return false
}

func (b *FastBoard) pointInDirection(m Move, cur uint16) Point {
	p := addPoints(IndexToPoint(cur, b.Width), moveToPoint(m))
	if b.IsWrapped {
		p = adjustForWrapped(p, b.Width, b.Height)
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

	newBoard.List = make([]Tile, len(b.List))
	copy(newBoard.List, b.List)

	newBoard.Ids = make(idsMap, len(b.Ids))
	for sId, id := range b.Ids {
		newBoard.Ids[sId] = id
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

	newBoard.hasEaten = make(map[SnakeId]bool, len(b.hasEaten))
	for id, he := range b.hasEaten {
		newBoard.hasEaten[id] = he
	}

	newBoard.Width = b.Width
	newBoard.Height = b.Height
	newBoard.IsWrapped = b.IsWrapped
	newBoard.hazardDamage = b.hazardDamage

	return newBoard
}

func (b *FastBoard) tileToString(index uint16) string {
	if b.IsTileSnakeHead(index) {
		return fmt.Sprintf(" %dh ", b.List[index].id)
	}
	if b.isTileSnakeSegment(index) {
		return fmt.Sprintf(" %ds ", b.List[index].id)
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
	for y := int(b.Height - 1); y >= 0; y-- {
		var line string
		for x := 0; x < int(b.Width); x++ {
			p := Point{X: int8(x), Y: int8(y)}
			line = line + b.tileToString(pointToIndex(p, b.Width))
		}
		fmt.Println(line)
	}
}

func (b *FastBoard) MoveToIndex(move SnakeMove) uint16 {
	p := pointInDirection(move.Dir, b.Heads[move.Id], b.Width, b.Height, b.IsWrapped)
	return pointToIndex(p, b.Width)
}
