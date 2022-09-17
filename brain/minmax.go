package brain

import (
	"fmt"
	"math"
	"time"

	g "github.com/nosnaws/tiam/game"
)

type mmNode struct {
	Score float64
	Move  g.Move
}

func Minmax(board *g.FastBoard, id g.SnakeId, depth int) g.Move {
	result := minmax(board, g.Left, id, depth, math.Inf(-1), math.Inf(1), true)

	return result.Move
}

func IdfsMinmax(b *g.FastBoard) g.Move {
	initial := Minmax(b, g.MeId, 2)

	duration, err := time.ParseDuration("150ms")
	if err != nil {
		panic("could not parse duration")
	}

	// create MAST table
	//mast := initMAST(*root.board)

	now := time.Now()
	currentDepth := 2
	bestMove := initial
miniloop:
	for timeout := time.After(duration); ; {
		select {
		case <-timeout:
			break miniloop
		default:
			currentDepth += 2
			fmt.Println("Running depth", currentDepth)
			bestMove = Minmax(b, g.MeId, currentDepth)
			fmt.Println("total time", time.Now().UnixMilli()-now.UnixMilli())
			fmt.Println("current best move", bestMove)
		}
	}

	fmt.Println("final time", time.Now().UnixMilli()-now.UnixMilli())

	fmt.Println("final move", bestMove)
	return bestMove
}

func minmax(board *g.FastBoard, move g.Move, id g.SnakeId, depth int, alpha, beta float64, maxPlayer bool) mmNode {
	if depth == 0 || board.IsGameOver() {
		return mmNode{
			Score: minmaxHeuristic(board, id),
			Move:  move,
		}
	}

	if maxPlayer {
		value := math.Inf(-1)
		curMove := move

		for _, m := range board.GetMovesForSnake(id) {
			//ns := board.Clone()
			//move := make(map[g.SnakeId]g.Move)
			//move[id] = m.Dir
			//ns.AdvanceBoardTurnBased(move)

			min := minmax(board, m.Dir, id, depth-1, alpha, beta, false)

			if min.Score > value {
				value = min.Score
				curMove = m.Dir
			}

			if value >= beta {
				break
			}

			if value > alpha {
				alpha = value
			}
		}
		return mmNode{Score: value, Move: curMove}

	} else {
		value := math.Inf(1)

		otherSnakes := getOtherSnakeIds(board, id)

		myMove := []g.SnakeMove{{Id: id, Dir: move}}
		movesTemp := [][]g.SnakeMove{myMove}
		for _, oId := range otherSnakes {
			moves := board.GetMovesForSnake(oId)

			//if len(movesTemp) < 1 {
			//movesTemp = append(movesTemp, moves)
			//} else {
			movesTemp = g.CartesianProduct(movesTemp, moves)
			//}
		}

		for _, moveSet := range movesTemp {
			moveMap := movesToMap(moveSet)
			ns := board.Clone()
			ns.AdvanceBoard(moveMap)

			max := minmax(&ns, move, id, depth-1, alpha, beta, true)

			if max.Score < value {
				value = max.Score
			}

			if value <= alpha {
				break
			}

			if value < beta {
				beta = value
			}
		}
		return mmNode{Score: value, Move: move}
	}
}

func getOtherSnakeIds(board *g.FastBoard, id g.SnakeId) []g.SnakeId {
	otherSnakes := []g.SnakeId{}
	for hid := range board.Heads {
		if board.IsSnakeAlive(hid) && hid != id && hid != g.MeId {
			otherSnakes = append(otherSnakes, hid)
		}
	}
	return otherSnakes
}

func minmaxHeuristic(board *g.FastBoard, id g.SnakeId) float64 {
	//health := float64(board.Healths[id])

	if board.IsGameOver() && board.IsSnakeAlive(id) {
		return math.MaxFloat64
	}

	if !board.IsSnakeAlive(id) {
		return -math.MaxFloat64
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
	numAlive := 0
	for id := range board.Lengths {
		if board.IsSnakeAlive(id) {
			numAlive += 1
		}
	}

	//total += health / 10

	total += float64(1000 / numAlive)
	//total += float64(board.Lengths[id]) / 2
	total += (100 / float64(voronoi.FoodDepth[id]))
	total += 8 * float64(voronoi.Score[id])

	return total // / math.Sqrt(100+math.Pow(total, 2))
}
