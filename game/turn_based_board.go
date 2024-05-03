package game

import (
	"math/rand"
)

// Advances the board, for use with a turn based board
func (b *FastBoard) AdvanceBoardTB(move SnakeMove) {
	id := move.Id
	dir := move.Dir
	//fmt.Println("STARTING MOVE", id, dir)
	//b.Print()
	//fmt.Println("move ", id, dir)
	oldHead := b.Heads[id]
	dirPoint := pointInDirection(dir, oldHead, b.Width, b.Height, b.isWrapped)

	// make sure we didn't move out of bounds
	// also prevents integer overflow
	if b.isOffBoard(dirPoint) {
		b.killTB(id)
		return
	}

	dirIndex := pointToIndex(dirPoint, b.Width)
	dirTile := b.List[dirIndex]

	if b.isTileSnakeSegment(dirIndex) {
		newSnakeId, _ := dirTile.GetSnakeId()
		newSnakeTail := b.getTBSnakeIdx(b.Heads[newSnakeId], newSnakeId)
		// head
		if dirTile.IsSnakeHead() {
			newL := b.Lengths[newSnakeId]
			curL := b.Lengths[id]

			if newL > curL {
				b.killTB(id)
				return
			}

			if newL == curL {
				b.killTB(id)
				b.killTB(newSnakeId)
				return
			}

			if newL < curL {
				b.killTB(newSnakeId)
			}
		} else if newSnakeTail != dirIndex || b.hasEaten[newSnakeId] {
			// body part
			b.killTB(id)
			return
		}
	}

	b.Healths[id] -= 1

	if dirTile.IsHazard() {
		b.Healths[id] -= b.hazardDamage * b.hazardDepth[dirIndex]
	}

	if dirTile.IsFood() {
		//fmt.Println("eating!")
		b.Healths[id] = 100
		b.Lengths[id] += 1
		b.hasEaten[id] = true

		tailIndex := b.getTBSnakeIdx(oldHead, id)
		if b.isTileSnakeHeadTail(tailIndex) {
			//fmt.Println("killing following snake")
			b.killTB(b.List[tailIndex].id)
		}
	} else {
		b.hasEaten[id] = false
	}

	tail := b.getTBSnakeIdx(oldHead, id)
	tailTile := b.List[tail]
	oldHeadTile := b.List[oldHead]

	//b.Print()
	//fmt.Println("tail index", tail)
	if oldHeadTile.IsTripleStack() {
		//fmt.Println("is triple stack")
		b.setTileSnakeDoubleStack(oldHead, id, dirIndex)
		b.setTBTileSnakeHead(dirIndex, id, oldHead)
		//fmt.Println("old head", oldHead)
	} else if tailTile.IsDoubleStack() {
		//fmt.Println("is double stack")
		//b.Print()
		b.setTileSnakeBodyPart(tail, id, oldHead)
		b.setTileSnakeBodyPart(oldHead, id, dirIndex)
		b.setTBTileSnakeHead(dirIndex, id, tail)
		//b.Print()
	} else if b.hasEaten[id] {
		//fmt.Println("has eaten")
		b.setTileSnakeBodyPart(oldHead, id, dirIndex)
		b.setTBTileSnakeHead(dirIndex, id, tail)
	} else {
		//fmt.Println("has not eaten, moving tail")
		b.setTileSnakeBodyPart(oldHead, id, dirIndex)
		b.setTBTileSnakeHead(dirIndex, id, b.getTBSnakeIdx(tail, id))

		// caveat for moving onto your own tail
		if tail != dirIndex {
			b.clearTileTB(tail, id)
		}
	}

	b.Heads[id] = dirIndex
	//fmt.Println("NEW BOARD STATE")
	//fmt.Println("Move was", id, dir)
	//b.Print()
}

func (b *FastBoard) getTBSnakeIdx(index uint16, id SnakeId) uint16 {
	if b.List[index].tailId == id {
		return b.List[index].tailIdx
	}
	return b.List[index].GetIdx()
}

func (b *FastBoard) setTBTileSnakeHead(index uint16, id SnakeId, nextIdx uint16) {
	// if this is another snakes tail
	if tailId, ok := b.getSnakeIdAtTile(index); ok && tailId != id {
		//fmt.Println("is other snake tail, setting head-tail")
		b.setTileSnakeHeadTail(index, id, nextIdx, tailId, b.List[index].GetIdx())
	} else {
		//fmt.Println("regular tile, clearing and setting head")
		b.clearTile(index)
		b.setTileSnakeHead(index, id, nextIdx)
	}
}

func (b *FastBoard) setTBTileSnakeTail(index uint16, id SnakeId, nextIdx uint16) {
	// if this is another snakes head
	if headId, ok := b.getSnakeIdAtTile(index); ok && headId != id {
		b.setTileSnakeHeadTail(index, headId, b.List[index].GetIdx(), id, nextIdx)
	} else {
		b.clearTile(index)
		b.setTileSnakeBodyPart(index, id, nextIdx)
	}
}

func (b *FastBoard) clearTileTB(index uint16, id SnakeId) {
	if b.isTileSnakeHeadTail(index) {
		if headId, _ := b.getSnakeIdAtTile(index); headId == id {
			// this is another snakes tail
			tile := b.List[index]
			b.clearTile(index)
			b.setTileSnakeBodyPart(index, tile.tailId, tile.tailIdx)
		} else {
			// this is another snakes head
			tile := b.List[index]
			b.clearTile(index)
			b.setTileSnakeHead(index, headId, tile.GetIdx())
		}
	} else {
		b.clearTile(index)
	}
}

// Kills a snake, for use with a turn based board
func (b *FastBoard) killTB(id SnakeId) {
	headIndex := b.Heads[id]
	tailIndex := b.getTBSnakeIdx(headIndex, id)
	//fmt.Println(id)
	//fmt.Println("head", headIndex)
	//fmt.Println("tail", tailIndex)
	//b.Print()
	//fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	seen := make(map[uint16]bool)

	if b.List[headIndex].IsTripleStack() {
		b.clearTileTB(headIndex, id)
	} else if b.List[tailIndex].IsDoubleStack() {
		b.clearTileTB(tailIndex, id)
		b.clearTileTB(headIndex, id)
	} else {
		currentIndex := tailIndex
		for currentIndex != headIndex {
			if seen[currentIndex] {
				b.Print()
				panic("in a loop!")
			}
			//fmt.Println(currentIndex)
			//fmt.Println(id)
			//b.Print()
			nextIndex := b.getTBSnakeIdx(currentIndex, id)

			b.clearTileTB(currentIndex, id)
			seen[currentIndex] = true
			currentIndex = nextIndex
		}

		b.clearTileTB(headIndex, id)
	}

	b.Lengths[id] = 0
	b.Healths[id] = 0
	b.Heads[id] = 0
	b.hasEaten[id] = false
}

// Runs a random game till completion, for use with a turn based board
func (b *FastBoard) RandomRolloutTB() {
	//turn := 0
	for !b.IsGameOver() {
		//moves := make([]SnakeMove, 0, len(b.Lengths))

		for id := range b.Lengths {
			if !b.IsSnakeAlive(id) {
				continue
			}

			var randomMove SnakeMove
			sMoves := b.GetMovesForSnakeTB(id)
			if len(sMoves) > 0 {
				randomMove = sMoves[rand.Intn(len(sMoves))]
			} else {
				randomMove = SnakeMove{id, Left}
			}

			//b.Print()
			//fmt.Printf("%s going %s", fmt.Sprint(id), randomMove.Dir)
			b.AdvanceBoardTB(SnakeMove{id, randomMove.Dir})
		}
		//turn += 1
	}
}

// Gets moves for a particular snake, for use with a turn based board
// Will return an empty array if there are no possible moves
func (b *FastBoard) GetMovesForSnakeTB(id SnakeId) []SnakeMove {
	var possibleMoves []SnakeMove
	snakeHeadIndex := b.Heads[id]

	moves := b.GetNeighborsTB(snakeHeadIndex)
	for _, m := range moves {
		possibleMoves = append(possibleMoves, SnakeMove{Id: id, Dir: m})
	}
	return possibleMoves
}

// Gets neighboring tiles, for use with a turn based board
// Does not return tiles that are off the board, insta kill hazards, and snake bodies
func (b *FastBoard) GetNeighborsTB(index uint16) []Move {
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

		if !b.isTileSnakeHead(dirIndex) && b.isTileNonHeadSnakeSegment(dirIndex) && !b.isTBTileSafeTail(dirIndex) {
			continue
		}

		possibleMoves = append(possibleMoves, dir)
	}

	return possibleMoves
}

func (b *FastBoard) isTBTileSafeTail(index uint16) bool {
	id, ok := b.getSnakeIdAtTile(index)
	if ok {
		headIndex := b.Heads[id]
		tailIndex := b.List[headIndex].GetIdx()

		return tailIndex == index && !b.hasEaten[id] && !b.List[tailIndex].IsDoubleStack()
	}

	return false
}
