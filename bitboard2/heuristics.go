package bitboard

import (
	"math"

	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/nosnaws/tiam/moveset"
)

// heuristic plan
// 1 / distance to closest food (consider finding path?)
// 1 / distance to closest smaller snake move (ditto)
// meVoronoi - sum of rest / all space  (this will be "area control")
// length
// maybe health?
// maybe some kind of center control bonus
// bonus for eliminating snakes?

type BasicStateWeights struct {
	Food   float64
	Aggr   float64
	Area   float64
	Length float64
	Health float64
	Tail   float64
}

// original scores
//total += foodScore * 2
//total += killScore
//total += areaScore * 0.8
//total += length * 0.5
//total += health * 0.1

// TODO: factor in depth for use during tree search
func (bb *BitBoard) BasicStateScore(snakeId string, depth int, weights BasicStateWeights) float64 {
	if !bb.IsSnakeAlive(snakeId) || moveset.IsEmpty(bb.GetMovesNoDefault(snakeId).Set) {
		return -100000 / (float64(depth) + 1)
	}

	if bb.IsGameOver() {
		return 100000 / (float64(depth) + 1)
	}

	snake := bb.GetSnake(snakeId)
	foodScore := bb.FoodScore(snakeId)
	killScore := bb.KillScore(snakeId)
	areaScore := bb.AreaControlScore(snakeId)
	tailPathScore := bb.TailPathScore(snakeId)

	length := float64(snake.Length)
	health := float64(snake.health)

	total := 0.0

	total += foodScore * weights.Food
	total += killScore * weights.Aggr
	total += areaScore * weights.Area
	total += length * weights.Length
	total += health * weights.Health
	total += tailPathScore * weights.Tail

	return total
}

func (bb *BitBoard) TailPathScore(snakeId string) float64 {
	score := 0.0
	snake := bb.GetSnake(snakeId)

	tailIdx := snake.getTailIndex()
	tailPathLength := BFS(bb, snake.GetHeadIndex(), func(i int) bool {
		if i == tailIdx {
			return true
		}
		return false
	})

	if tailPathLength != -1 {
		score = 1
	}

	return score
}

func (bb *BitBoard) KillScore(snakeId string) float64 {
	smallerSnakes := bb.getSmallerSnakes(snakeId)

	distances := []float64{}
	for _, ss := range smallerSnakes {
		distances = append(distances, snakeMoveDistances(bb, snakeId, ss)...)
	}

	return 1 / (1 + min(distances))
}

func (bb *BitBoard) FoodScore(snakeId string) float64 {
	distances := foodDistances(bb, snakeId)

	return 1 / (1 + min(distances))
}

func (bb *BitBoard) AreaControlScore(snakeId string) float64 {
	voronoi := bb.Voronoi()

	score := float64(voronoi.Score[snakeId])
	enemyTotal := 0

	for id, s := range voronoi.Score {
		if id != snakeId {
			enemyTotal += s
		}
	}

	return (score - float64(enemyTotal)) / (float64(bb.width) * float64(bb.height))
}

func manhattanDistance(p1, p2 api.Coord) float64 {
	return math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
}

func foodDistances(b *BitBoard, snakeId string) []float64 {
	distances := []float64{}
	boardLen := b.width * b.height
	snake := b.GetSnake(snakeId)
	headP := indexToPoint(snake.GetHeadIndex(), b.width)

	for i := 0; i < boardLen; i++ {
		if b.IsIndexFood(i) {
			point := indexToPoint(i, b.width)
			distances = append(distances, manhattanDistance(headP, point))
		}
	}

	return distances
}

func snakeMoveDistances(b *BitBoard, meId string, snakeId string) []float64 {
	me := b.GetSnake(meId)
	snake := b.GetSnake(snakeId)
	meHead := me.GetHeadIndex()
	snakeHead := snake.GetHeadIndex()
	mePoint := indexToPoint(meHead, b.width)

	snakeMoves := moveset.Split(b.GetMoves(snakeId).Set)

	distances := []float64{}
	for _, m := range snakeMoves {
		mIdx := indexInDirection(MoveSetToDir(m), snakeHead, b.width, b.height, b.isWrapped)
		mPoint := indexToPoint(mIdx, b.width)
		distances = append(distances, manhattanDistance(mePoint, mPoint))
	}

	return distances
}

func sum(args []float64) float64 {
	total := 0.0

	for _, v := range args {
		total += v
	}

	return total
}

func max(args []float64) float64 {
	if len(args) == 0 {
		return 0
	}
	m := args[0]

	for _, v := range args {
		if v > m {
			m = v
		}
	}
	return m
}

func min(args []float64) float64 {
	if len(args) == 0 {
		return 0
	}
	m := args[0]

	for _, v := range args {
		if v < m {
			m = v
		}
	}
	return m
}

func (bb *BitBoard) getSmallerSnakes(snakeId string) []string {
	smaller := []string{}
	length := bb.GetSnake(snakeId).Length
	for id, snake := range bb.Snakes {
		if id != snakeId && snake.Length < length {
			smaller = append(smaller, id)
		}
	}

	return smaller
}

func (bb *BitBoard) getLargerSnakes(snakeId string) []string {
	smaller := []string{}
	length := bb.GetSnake(snakeId).Length
	for id, snake := range bb.Snakes {
		if id != snakeId && snake.Length >= length {
			smaller = append(smaller, id)
		}
	}

	return smaller
}
