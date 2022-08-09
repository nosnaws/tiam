package brain

import (
	"math"

	g "github.com/nosnaws/tiam/game"
)

type abReturn struct {
	value int32
	Move  g.SnakeMove
}

func AlphaBeta(state g.FastBoard, depth int32, alpha int32, beta int32, maxPlayer bool, maxPlayerId g.SnakeId, lastMove g.SnakeMove) abReturn {
	if depth == 0 || state.IsGameOver() {
		return abReturn{value: evalState(state, maxPlayerId), Move: lastMove}
	}

	if maxPlayer {
		value := abReturn{value: int32(math.Inf(-1)), Move: lastMove}

		for _, move := range state.GetMovesForSnake(maxPlayerId) {
			mMap := make(map[g.SnakeId]g.Move, 1)
			mMap[move.Id] = move.Dir
			ns := state.Clone()

			ns.AdvanceBoard(mMap)

			result := AlphaBeta(ns, depth-1, alpha, beta, false, maxPlayerId, move)
			value = takeMax(value, result)

			if value.value >= beta {
				break
			}
			alpha = maxInt(alpha, value.value)

		}
		return value
	} else {
		value := abReturn{value: int32(math.Inf(1)), Move: lastMove}

		var otherMoveSet [][]g.SnakeMove
		for id := range state.Lengths {
			if id != maxPlayerId {
				otherMoveSet = append(otherMoveSet, state.GetMovesForSnake(id))
			}
		}
		maxPlayerMove := []g.SnakeMove{lastMove}
		possibleStates := g.CartesianProduct(otherMoveSet, maxPlayerMove)

		for _, moves := range possibleStates {
			mMap := make(map[g.SnakeId]g.Move, 1)
			for _, m := range moves {
				if m.Id != maxPlayerId {
					mMap[m.Id] = m.Dir
				}
			}

			ns := state.Clone()
			ns.AdvanceBoard(mMap)

			result := AlphaBeta(ns, depth-1, alpha, beta, true, maxPlayerId, lastMove)
			value = takeMin(value, result)

			if value.value <= alpha {
				break
			}

			beta = minInt(beta, value.value)
		}
		return value
	}
}

func evalState(state g.FastBoard, maxPlayerId g.SnakeId) int32 {
	if state.IsGameOver() {
		if state.IsSnakeAlive(maxPlayerId) {
			return math.MaxInt32
		}
	}

	if !state.IsSnakeAlive(maxPlayerId) {
		return -math.MaxInt32
	}

	total := int32(0)
	total += int32(state.Healths[maxPlayerId] / 100)
	total += int32(state.Lengths[maxPlayerId])

	return total
}

func maxInt(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func takeMax(a, b abReturn) abReturn {
	if a.value > b.value {
		return a
	}
	return b
}

func takeMin(a, b abReturn) abReturn {
	if a.value < b.value {
		return a
	}
	return b
}
