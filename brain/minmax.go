package brain

import (
	"math"

	g "github.com/nosnaws/tiam/game"
)

type mmNode struct {
	Score float64
	Move  g.Move
}

func Minmax(board *g.FastBoard, move g.Move, depth int, alpha, beta float64, maxPlayer bool) mmNode {
	if depth == 0 || board.IsGameOver() {
		return mmNode{
			Score: minmaxHeuristic(board, g.MeId),
			Move:  move,
		}
	}

	if maxPlayer {
		value := -math.MaxFloat64
		curMove := move

		for _, m := range board.GetMovesForSnake(g.MeId) {
			ns := board.Clone()
			move := make(map[g.SnakeId]g.Move)
			move[g.MeId] = m.Dir
			ns.AdvanceBoardTurnBased(move)

			result := Minmax(&ns, m.Dir, depth-1, alpha, beta, false)

			if result.Score > value {
				value = result.Score
				curMove = result.Move
			}

			if value >= beta {
				break
			}
		}
		return mmNode{Score: value, Move: curMove}

	} else {
		value := math.MaxFloat64

		otherSnakes := getOtherSnakeIds(board)

		movesTemp := [][]g.SnakeMove{}
		for _, id := range otherSnakes {
			moves := board.GetMovesForSnake(id)

			if len(movesTemp) < 1 {
				movesTemp = append(movesTemp, moves)
			} else {
				movesTemp = g.CartesianProduct(movesTemp, moves)
			}
		}

		for _, moveSet := range movesTemp {
			moveMap := movesToMap(moveSet)
			ns := board.Clone()
			ns.AdvanceBoardTurnBased(moveMap)

			result := Minmax(&ns, move, depth-1, alpha, beta, true)

			if value < result.Score {
				value = result.Score
			}

			if value <= alpha {
				break
			}
		}
		return mmNode{Score: value, Move: move}
	}
}

func getOtherSnakeIds(board *g.FastBoard) []g.SnakeId {
	otherSnakes := []g.SnakeId{}
	for id := range board.Heads {
		if board.IsSnakeAlive(id) && id != g.MeId {
			otherSnakes = append(otherSnakes, id)
		}
	}
	return otherSnakes
}

func minmaxHeuristic(board *g.FastBoard, id g.SnakeId) float64 {
	health := float64(board.Healths[id])

	if !board.IsSnakeAlive(id) {
		return -1.0
	}

	//isLargestSnake := true
	//for sId, l := range node.board.Lengths {
	//if sId != id && l >= node.board.Lengths[id] {
	//isLargestSnake = false
	//}
	//}

	voronoi := g.Voronoi(board, id)

	total := 0.0
	//if isLargestSnake {
	//total += config.BigSnakeReward
	//}

	total += 1.2 * ((health - float64(voronoi.FoodDepth)) / 14.5)
	total += 1 * float64(voronoi.Score)

	return total / math.Sqrt(100+math.Pow(total, 2))
}
