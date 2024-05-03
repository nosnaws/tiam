package mcts

import (
	"math/rand"
	"time"

	"github.com/nosnaws/tiam/board"
)

func StrategicRollout(fb *board.FastBoard) {
	rand.Seed(time.Now().UnixMilli())
	moves := make(map[board.SnakeId]board.Move, len(fb.Lengths))
	for !fb.IsGameOver() {
		for id := range fb.Lengths {
			if !fb.IsSnakeAlive(id) {
				moves[id] = ""
				continue
			}

			move := mixedStrategy(fb, id)
			moves[id] = move.Dir
		}
		fb.AdvanceBoard(moves)
	}
}

func RandomRollout(fb *board.FastBoard) {
	rand.Seed(time.Now().UnixMilli())
	moves := make(map[board.SnakeId]board.Move, len(fb.Lengths))
	for !fb.IsGameOver() {
		for id := range fb.Lengths {
			if !fb.IsSnakeAlive(id) {
				moves[id] = ""
				continue
			}

			randomMove := randomStrategy(fb, id)
			moves[id] = randomMove.Dir
		}
		fb.AdvanceBoard(moves)
	}
}
