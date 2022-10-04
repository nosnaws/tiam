package mmm

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	b "github.com/nosnaws/tiam/board"
)

type mmNode struct {
	Score float64
	Move  b.Move
}

func MultiMinmaxID(board *b.FastBoard, cache *Cache) b.Move {
	currentDepth := 0

	// current issue, if it gets to a deep enough depth, it won't return in time
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*350)
	defer cancel()

	//startTime := time.Now().UnixMilli()

	lastBestMove := b.Move("")
	//lastBestScore := math.Inf(-1)

	if cache == nil {
		cache = CreateCache(board, currentDepth)
	}

mmloop:
	for {
		select {
		case <-ctx.Done():
			break mmloop
		default:
			currentDepth += 2
			cache.setCurMax(currentDepth)

			if currentDepth > 500 {
				fmt.Println("CUTTING OFF")
				break mmloop
			}

			fmt.Println("RUNNING ITERATION", currentDepth)
			newMove, newScore, ignoreScore := MultiMinmax(ctx, cache, board, currentDepth, lastBestMove)
			fmt.Println("ITERATION RESULTS", newMove, newScore)

			if !ignoreScore {
				fmt.Println("UPDATING BEST")
				lastBestMove = newMove
				//lastBestScore = newScore
			}
		}
	}

	return lastBestMove
}

func MultiMinmax(ctx context.Context, cache *Cache, board *b.FastBoard, depth int, firstMove b.Move) (b.Move, float64, bool) {
	rand.Seed(time.Now().UnixMicro())

	depthLogger := createLogger(depth)
	depthLogger.turnLoggerOff()
	selectedMove := b.Left
	maxScore := math.Inf(-1)
	maxOpp := b.SnakeId(0)
	alpha := math.Inf(-1)

	moves := b.Shuffle(board.GetMovesForSnake(b.MeId))
	if firstMove != "" {
		sort.Slice(moves, func(i, j int) bool {
			return moves[i].Dir == firstMove
		})
	}

	ignoreScore := false
	for _, maxMove := range board.GetMovesForSnake(b.MeId) {
		ns := board.Clone()
		var minScore float64
		var lastOpp b.SnakeId
		minScore, lastOpp, ignoreScore = multiHelper(ctx, cache, &ns, maxMove, depth, alpha, maxOpp)

		depthLogger.debug(depth, "maxScore: ", maxScore)
		depthLogger.debug(depth, "minScore: ", minScore)
		depthLogger.debug(depth, "currentMove: ", selectedMove)
		if minScore >= maxScore {
			depthLogger.debug(depth, "min greater than max, using new move", maxMove.Dir)
			selectedMove = maxMove.Dir
			maxScore = minScore
			maxOpp = lastOpp
			alpha = maxScore
		}
	}

	return selectedMove, maxScore, ignoreScore
}

func multiHelper(ctx context.Context, cache *Cache, board *b.FastBoard, maxMove b.SnakeMove, depth int, alpha float64, firstOpp b.SnakeId) (float64, b.SnakeId, bool) {
	depthLogger := getLogger()
	minScore := math.Inf(1)
	beta := math.Inf(1)
	otherSnakes := getOtherSnakeIds(board, b.MeId)

	ignoreResults := false
	oSnakes := b.Shuffle(otherSnakes)
	if firstOpp != 0 {
		sort.Slice(oSnakes, func(i, j int) bool {
			return oSnakes[i] == firstOpp
		})
	}

	minOpp := b.SnakeId(0)
	for _, oSId := range oSnakes {
		ns := board.Clone()
		//removeOtherSnakes(&ns, b.MeId, oSId)

		depthLogger.debug(depth, fmt.Sprintf("%d:%s enemyId:%d", b.MeId, maxMove.Dir, oSId))
		var min float64
		min, ignoreResults = minmax(ctx, cache, &ns, maxMove, b.MeId, oSId, depth-1, alpha, beta, false)
		depthLogger.debug(depth, fmt.Sprintf("%d:%s enemyId:%d - score:%f", b.MeId, maxMove.Dir, oSId, min))
		if min < minScore {
			depthLogger.debug(depth, "result smaller than min, updating min")
			minScore = min
			minOpp = oSId
		}

		beta = minScore

		depthLogger.debug(depth, "alpha: ", alpha)
		depthLogger.debug(depth, "beta: ", beta)
		if alpha >= beta {
			fmt.Println("alpha greater than beta, breaking")
			break
		}
	}

	return minScore, minOpp, ignoreResults
}

func minmax(ctx context.Context, cache *Cache, board *b.FastBoard, maxMove b.SnakeMove, maxId b.SnakeId, minId b.SnakeId, depth int, alpha, beta float64, maxPlayer bool) (float64, bool) {
	ignoreResults := false
	select {
	case <-ctx.Done():
		ignoreResults = true
		fmt.Println("ALL DONE, EXITING", depth)
		return 0, ignoreResults
	default:
	}
	log := getLogger()
	if depth == 0 || shouldExit(board, maxId, minId) {
		return minmaxHeuristic(board, maxId, minId, depth), ignoreResults
	}

	// transposition table logic
	if entry, ok := cache.getEntry(board, minId, depth); ok {
		if entry.isExact() {
			return entry.value, ignoreResults
		}

		if entry.isLowerBound() && entry.value > alpha {
			alpha = entry.value
		}
		if entry.isUpperBound() && entry.value < beta {
			beta = entry.value
		}

	}

	// issues with heuristic
	// since one snake never moves, the game might never be considered "over"
	// fix: added extra logic to check if they have no moves, which would be considered dead
	// since there are only 2 players here, this should work fine

	// issues with no moves
	// need to cut out when the snake has no moves
	// also an issue with continuing the recursion when snakes are dead

	var moveScore float64
	if maxPlayer {
		moveScore = math.Inf(-1)

		maxMoves := b.Shuffle(board.GetMovesForSnakeNoDefault(maxId))
		//maxMoves := b.Shuffle(getMovesForSnakePruned(board, maxId)) // appears to be too much pruning
		if len(maxMoves) == 0 {
			return minmaxHeuristic(board, maxId, minId, depth), ignoreResults
		}

		for _, m := range maxMoves {
			log.debug(depth, "considering MAX move", m.Dir, minId)
			var result float64
			result, ignoreResults = minmax(ctx, cache, board, m, maxId, minId, depth-1, alpha, beta, false)
			log.debug(depth, "result", result)

			if result > moveScore {
				moveScore = result
			}

			if moveScore > alpha {
				alpha = moveScore
			}

			if alpha >= beta {
				break
			}
		}
	} else {
		moveScore = math.Inf(1)

		minMoves := b.Shuffle(board.GetMovesForSnakeNoDefault(minId))
		//minMoves := b.Shuffle(getMovesForSnakePruned(board, minId))
		if len(minMoves) == 0 {
			return minmaxHeuristic(board, maxId, minId, depth), ignoreResults
		}

		for _, minMove := range minMoves {
			moveMap := b.MovesToMap([]b.SnakeMove{maxMove, minMove})
			ns := board.Clone()
			ns.AdvanceBoard(moveMap)

			log.debug(depth, "considering MIN move", minMove.Dir, minId)
			var result float64
			result, ignoreResults = minmax(ctx, cache, &ns, maxMove, maxId, minId, depth-1, alpha, beta, true)
			log.debug(depth, "result", result)

			if result < moveScore {
				moveScore = result
			}

			if moveScore < alpha {
				beta = moveScore
			}

			if alpha >= beta {
				break
			}
		}
	}

	// adding entries to transposition table
	if moveScore <= alpha {
		cache.addUpperBound(board, moveScore, minId, depth)
	}
	if moveScore > alpha && moveScore < beta {
		cache.addExact(board, moveScore, minId, depth)
	}
	if moveScore >= beta {
		cache.addLowerBound(board, moveScore, minId, depth)
	}

	return moveScore, ignoreResults
}

func shouldExit(board *b.FastBoard, maxId b.SnakeId, minId b.SnakeId) bool {
	if board.IsGameOver() {
		return true
	}

	if isDeadOrOut(board, maxId) || isDeadOrOut(board, minId) {
		return true
	}

	return false
}

func isDeadOrOut(board *b.FastBoard, id b.SnakeId) bool {
	if !board.IsSnakeAlive(id) {
		return true
	}

	moves := board.GetMovesForSnakeNoDefault(id)
	if len(moves) == 0 {
		return true
	}

	return false
}

func getOtherSnakeIds(board *b.FastBoard, id b.SnakeId) []b.SnakeId {
	otherSnakes := []b.SnakeId{}
	for hid := range board.Heads {
		if board.IsSnakeAlive(hid) && hid != id {
			otherSnakes = append(otherSnakes, hid)
		}
	}
	return otherSnakes
}

func removeOtherSnakes(board *b.FastBoard, id1 b.SnakeId, id2 b.SnakeId) {
	for hid := range board.Heads {
		if hid != id1 && hid != id2 {
			board.RemoveSnake(hid)
		}
	}
}

func getMovesForSnakePruned(board *b.FastBoard, id b.SnakeId) []b.SnakeMove {
	moves := board.GetMovesForSnakeNoDefault(id)

	prunedMoves := []b.SnakeMove{}
	for _, m := range moves {
		if len(board.GetNeighbors(board.IndexInDirection(m.Dir, board.Heads[id]))) > 0 {
			prunedMoves = append(prunedMoves, m)
		}
	}

	if len(prunedMoves) == 0 {
		return moves
	}

	return prunedMoves
}

func minmaxHeuristic(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	//health := float64(board.Healths[id])
	log := getLogger()

	if isDeadOrOut(board, maxId) {
		log.debugAlways(depth, "We lose!")
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		log.debugAlways(depth, "We win!")
		return float64(10000 * (depth + 1))
	}

	minMoves := board.GetMovesForSnakeNoDefault(minId)
	minMIndexs := []uint16{}
	for _, mm := range minMoves {
		minMIndexs = append(minMIndexs, b.IndexInDirection(mm.Dir, board.Heads[minId], board.Width, board.Height, board.IsWrapped))
	}
	//voronoi := b.Voronoi(board, id)
	maxTail := board.GetSnakeTail(maxId)
	ff, foodDepth, tailDepth := floodfill(board, int(board.Heads[maxId]), int(board.Lengths[maxId]*2), []uint16{maxTail})

	total := 0.0
	//total += float64(board.Lengths[maxId]) * 0.02
	//total += float64(voronoi.Score[id]) * 0.01
	//total += float64(voronoi.FoodDepth[id]) * 0.01
	total += float64(ff) / 2

	if ff < int(board.Lengths[maxId]) {
		// less space available than our length is bad
		// but not quite as bad as losing the game
		total += float64(-50 * (depth + 1))
	}

	if board.Lengths[maxId] > board.Lengths[minId] {
		//total += float64(-enemyDepth)
	} else if foodDepth != -1 {
		total += float64(-foodDepth)
	}
	// if we are near a hazard, that's not the best
	if hNeighbors := board.GetHazardNeighbors(board.Heads[maxId]); len(hNeighbors) > 0 {
		total += float64(-20 * len(hNeighbors))
	}

	// moving toward our tail is usually good
	if tailDepth != -1 {
		total += float64(-tailDepth * 2)
	} else {
		total += float64(-50 * (depth + 1))
	}

	return total
}

// TODO:
//func MultiMinmaxAllCylinders(board *b.FastBoard, depth int) b.Move {
//depthLogger := createLogger(depth)
//selectedMove := b.Left
//maxScore := math.Inf(-1)
//alpha := math.Inf(-1)

//for _, maxMove := range board.GetMovesForSnake(b.MeId) {
//ns := board.Clone()
//minScore := multiHelper(&ns, maxMove, depth, alpha)

//depthLogger.debug(depth, "maxScore: ", maxScore)
//depthLogger.debug(depth, "minScore: ", minScore)
//depthLogger.debug(depth, "currentMove: ", selectedMove)
//if minScore >= maxScore {
//depthLogger.debug(depth, "min greater than max, using new move", maxMove.Dir)
//selectedMove = maxMove.Dir
//maxScore = minScore
//alpha = maxScore
//}
//}

//board.Print()
//return selectedMove
//}
