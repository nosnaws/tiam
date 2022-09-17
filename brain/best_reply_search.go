package brain

import (
	"fmt"
	"log"
	"math"
	"sort"
	"sync"

	g "github.com/nosnaws/tiam/game"
	"golang.org/x/net/context"
)

type BRSNode struct {
	Score float64
	Move  g.Move
}

type bestMove struct {
	m        BRSNode
	exitFlag bool
	mutex    sync.RWMutex
}

type BRSGame struct {
	//exitFlag bool
	//mutex    sync.RWMutex
	cache map[BoardHash]ttEntry
	tt    TranspositionTable
}

func CreateBRSGame(b *g.FastBoard) *BRSGame {
	brsGame := BRSGame{}
	brsGame.cache = createTTTable()
	brsGame.tt = InitializeTranspositionTable(int(b.Height), int(b.Width))
	return &brsGame
}

func IDBRS(ctx context.Context, b *g.FastBoard) g.Move {
	brsGame := CreateBRSGame(b)

	var finalMove g.Move
	var brsDone chan BRSNode
	runBrs := make(chan int, 1)
	currentDepth := 1

	for {

		// time to run another depth
		if brsDone == nil {
			currentDepth += 2
			runBrs <- currentDepth
		}

		select {
		case depth := <-runBrs:
			brsDone = make(chan BRSNode)
			go func() {
				log.Print("STARTING DEPTH ", depth)
				result := brsGame.BRS(b, g.Left, depth, math.Inf(-1), math.Inf(1), true)
				select {
				case brsDone <- result:
					log.Print("SENDING RESULT ", result.Move, depth)
				case <-ctx.Done():
					log.Print("EXITING ", depth)
				}
			}()
		case res := <-brsDone:
			brsDone = nil
			finalMove = res.Move
			log.Print("updating best move ", finalMove)
			if currentDepth > 50 {
				log.Println("over 50 depth, exiting")
				return finalMove
			}
		case <-ctx.Done():
			log.Println("returning final move ", finalMove)
			return finalMove
		}
	}
}

func (brs *BRSGame) BRS(board *g.FastBoard, move g.Move, depth int, alpha, beta float64, maxPlayer bool) BRSNode {
	//ogAlpha := alpha
	//bHash := HashBoard(brs.tt, *board)
	//if entry, ok := brs.cache[bHash]; ok && entry.depth >= depth {
	////fmt.Println("cache hit!")
	//if entry.flag == EXACT {
	//return entry.value
	//} else if entry.flag == LOWERBOUND {
	//if entry.value.Score > alpha {
	//alpha = entry.value.Score
	//}
	//} else if entry.flag == UPPERBOUND {
	//if entry.value.Score < beta {
	//beta = entry.value.Score
	//}
	//}

	//if alpha >= beta {
	//return entry.value
	//}
	//}

	if depth <= 0 || board.IsGameOver() || !board.IsSnakeAlive(g.MeId) {
		return BRSNode{
			Score: brsHeuristic(board, depth),
			Move:  move,
		}
	}

	moves := []g.SnakeMove{}
	if maxPlayer {
		orderedMoves := board.GetMovesForSnakeTB(g.MeId)

		for _, m := range orderedMoves {
			print(depth, "adding move", m, depth)
			moves = append(moves, m)
		}
		//if len(moves) < 1 {
		//moves = append(moves, g.SnakeMove{Id: g.MeId, Dir: g.Left})
		//}
		//fmt.Println("My turn moves", moves)
	} else {
		otherSnakes := getOtherSnakeIds(board, g.MeId)
		for _, id := range otherSnakes {
			if !board.IsSnakeAlive(id) {
				continue
			}
			for _, m := range board.GetMovesForSnakeTB(id) {
				print(depth, "adding enemy move", m, depth)
				moves = append(moves, m)
			}
			if len(moves) < 1 {
				moves = append(moves, g.SnakeMove{Id: id, Dir: g.Left})
			}
		}
		//fmt.Println("Opponents turn moves", moves)
	}

	// short curcuit if we dont' have any moves
	if maxPlayer && len(moves) == 0 {
		ns := board.Clone()
		// go left so the game can end
		ns.AdvanceBoard(movesToMap([]g.SnakeMove{{Id: g.MeId, Dir: g.Left}}))

		return BRSNode{
			Score: brsHeuristic(&ns, depth),
			Move:  move,
		}
	}

	bestMove := move
	if maxPlayer {
		bestMove = moves[0].Dir
	}

	for _, m := range moves {
		ns := board.Clone()

		myMove := move
		if maxPlayer {
			myMove = m.Dir
		}

		moves := []g.SnakeMove{m}

		if !maxPlayer {
			otherSnakeIds := getOtherSnakeIds(board, m.Id)
			for _, oId := range otherSnakeIds {
				moves = append(moves, g.SnakeMove{Id: oId, Dir: g.Left})
			}
			// run moves in static order to prevent weird states
			sort.Slice(moves, func(i, j int) bool {
				return moves[i].Id < moves[j].Id
			})
		}

		for _, move := range moves {
			ns.AdvanceBoardTB(move)
		}
		//ns.AdvanceBoardTB(m)

		result := brs.BRS(&ns, myMove, depth-1, -beta, -alpha, !maxPlayer)
		rScore := -result.Score
		print(depth, "result", rScore, myMove, depth)

		if rScore >= beta {
			print(depth, "result score greater than beta, returning", rScore, myMove, maxPlayer, depth)
			return BRSNode{Score: rScore, Move: myMove}
		}
		if rScore > alpha {
			print(depth, "Updating alpha", rScore, depth)
			alpha = rScore
			if maxPlayer {
				print(depth, "max player, updating move", result.Move, depth)
				bestMove = result.Move
			}
		}
	}

	//newEntry := ttEntry{}
	//newEntry.value = BRSNode{Score: alpha, Move: bestMove}
	//if alpha <= ogAlpha {
	//newEntry.flag = UPPERBOUND
	//} else if alpha >= beta {
	//newEntry.flag = LOWERBOUND
	//} else {
	//newEntry.flag = EXACT
	//}
	//newEntry.depth = depth
	//brs.cache[bHash] = newEntry

	print(depth, "out of loop returning", alpha, bestMove, maxPlayer, depth)
	return BRSNode{Score: alpha, Move: bestMove}
}

func print(depth int, s ...any) {
	if depth <= 5 {
		fmt.Println(s...)
	}
}

func brsHeuristic(board *g.FastBoard, depth int) float64 {
	//board.Print()
	id := g.MeId
	//health := float64(board.Healths[id])

	if !board.IsSnakeAlive(id) {
		fmt.Println("game over we lose")
		return -10000 * float64(depth+1)
	}

	if board.IsGameOver() {
		fmt.Println("game over we win")
		return 10000 * float64(depth+1)
	}

	//if !board.IsSnakeAlive(id) {
	//return -math.MaxFloat64
	//}

	//isLargestSnake := true
	//for sId, l := range node.board.Lengths {
	//if sId != id && l >= node.board.Lengths[id] {
	//isLargestSnake = false
	//}
	//}

	voronoi := g.Voronoi(board, id)
	//fmt.Println("voronoi", voronoi.Score)
	//fmt.Println("food", voronoi.FoodDepth)

	total := 0.0
	//if isLargestSnake {
	//total += config.BigSnakeReward
	//}
	//numAlive := 0
	//for id := range board.Lengths {
	//if board.IsSnakeAlive(id) {
	//numAlive += 1
	//}
	//}

	//total += 10 * (101 / (health + 1))
	total += 2 * float64(board.Lengths[id])

	//total += float64(100 / numAlive)
	//total += float64(board.Lengths[id]) / 2
	//total += 2 * float64(board.Lengths[id])
	if voronoi.FoodDepth[id] > -1 {
		total += 8 * (1 / float64(voronoi.FoodDepth[id]+1))
	}
	total += 12 * float64(voronoi.Score[id])

	return total // / math.Sqrt(100+math.Pow(total, 2))
}

const (
	UPPERBOUND = "upper"
	LOWERBOUND = "lower"
	EXACT      = "exact"
)

type ttEntry struct {
	value BRSNode
	flag  string
	depth int
}

func createTTTable() map[BoardHash]ttEntry {
	return make(map[BoardHash]ttEntry)
}
