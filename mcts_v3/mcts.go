package mctsv3

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	api "github.com/nosnaws/tiam/battlesnake"
	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/moveset"
)

// TODO: the bitboard is going to remove snakes when they die,
// this does not mesh well with the way we have rewards structured and backpropagation.
// To avoid having to figure out a different way of doing things here, I need to keep track
// of the snakes present in the game in a long living (as long as a game) reference so that MCTS can loop through all the snakes
//

// 0 - survival, 1 - food, 2 - aggressive
type tactics [3]float64

const survivalTac = tactic(0)
const foodTac = tactic(1)
const aggTac = tactic(2)

type tactic int
type reward struct {
	tac   tactics
	plays float64
}
type rewards map[moveset.MoveSet]reward
type rewardMatrix map[string]rewards
type rolloutScore map[string]tactics

type MCTSConfig struct {
	MinimumSims    int
	SimLimit       int
	TreeDepthLimit int
}

type node struct {
	b             *bitboard.BitBoard
	previousMoves map[string]moveset.MoveSet
	prevMoves     []bitboard.SnakeMoveSet
	depth         int
	parent        *node
	plays         float64
	children      []*node
	reward        rewardMatrix
	me            string
	gtValue       float64 // Game-theoretical value https://www.cs.drexel.edu/~santi/teaching/2012/CS680/papers/W4R2.pdf
}

var globalRoot *node
var globalZobrist *bitboard.ZobristTable

func ChooseNewRoot(root *node, state api.GameState) *node {
	//return CreateTree(bitboard.CreateBitBoard(state), state.You.ID)

	if root == nil {
		return CreateTree(bitboard.CreateBitBoard(state), state.You.ID)
	} //fmt.Println("ROOT")
	//printNode(globalRoot)

	newState := root.b.Clone()
	newRoot := findMatchingChild(root, state)

	if newRoot != nil {
		newState.AdvanceTurn(newRoot.prevMoves)
		newState.Update(bitboard.CreateBitBoard(state))
		newRoot.b = newState
		newRoot.parent = nil
		newRoot.prevMoves = []bitboard.SnakeMoveSet{}
		newRoot.previousMoves = make(map[string]moveset.MoveSet)
		decayTree(root, 0.6)

		fmt.Println("SELECTED")
		printNode(newRoot, state.You.ID)
		return newRoot
	}

	fmt.Println("DROPPING TREE")
	return CreateTree(bitboard.CreateBitBoard(state), state.You.ID)
}

func printNode(n *node, snakeId string) {
	if n.b != nil {
		n.b.Print()
	}
	for m, r := range n.reward[snakeId] {
		fmt.Println("CHILD", r.plays, m, r.tac)
	}
}

// There is currently a problem where the state gets messed up and then this can't properly find the correct child
// It appears to have something to do with snakes not being removed from the game when they die.
// This could be because we selected the wrong child at an earlier turn, but I'm not sure.
// I noticed this because the number of children was larger than it should be. The test game had killed a snake, but it was still alive and well within the game tree
func findMatchingChild(root *node, state api.GameState) *node {
	actualBoard := bitboard.CreateBitBoard(state)

	for _, child := range root.children {
		ns := root.b.Clone()
		ns.AdvanceTurn(child.prevMoves)
		if ns.IsEqual(actualBoard) {
			return child
		}
	}
	return nil

	//actualMoves := []bitboard.SnakeMove{}
	//for _, snake := range state.Board.Snakes {
	//actualMoves = append(actualMoves, bitboard.SnakeMove{
	//Id:  snake.ID,
	//Dir: root.b.GetLastSnakeMoveFromExternal(snake),
	//})
	//}

	//for _, child := range root.children {
	//fmt.Println("CHILD PLAYS", child.prevMoves, child.reward)
	//}

	//for _, child := range root.children {
	//if isStateEqual(actualMoves, child.previousMoves) {
	//return child
	//}
	//}

	//return nil
}

func updateDepths(n *node) {
	if n == nil {
		return
	}
	n.depth -= 1
	for _, c := range n.children {
		updateDepths(c)
	}
}

func decayTree(n *node, decayConst float64) {
	if n == nil {
		return
	}
	n.depth -= 1

	for id, re := range n.reward {
		snakeReward := make(rewards)
		for m, mRe := range re {
			newReward := reward{}
			decayed := tactics{0, 0, 0}

			for i := 0; i < len(mRe.tac); i++ {
				decayed[i] = mRe.tac[i] * decayConst
			}

			newReward.tac = decayed
			newReward.plays = mRe.plays * decayConst
			snakeReward[m] = newReward
		}
		n.reward[id] = snakeReward
	}
	n.plays = n.plays * decayConst

	for _, c := range n.children {
		decayTree(c, decayConst)
	}
}

func CreateTree(b *bitboard.BitBoard, me string) *node {
	r := createNode(b, []bitboard.SnakeMoveSet{}, nil, me)
	//globalRoot = r
	//globalRoot.b = b
	r.b = b
	return r
}

func MCTS(root *node, minSims, simLimit, depthLim int) bitboard.Dir {
	duration, _ := time.ParseDuration("400ms")
	t := determineTactic(root)
	//availableWorkers := runtime.GOMAXPROCS(0)
	availableWorkers := 1
	fmt.Println("WORKERS", availableWorkers)
	rand := rand.New(rand.NewSource(time.Now().Unix()))

	maxDepth := 0
mctsloop:
	for timeout := time.After(duration); ; {
		select {
		case <-timeout:
			break mctsloop
		default:
			state := root.b.Clone()
			node := selectNode(root, state, minSims, depthLim, t)
			if node.depth > maxDepth {
				maxDepth = node.depth
			}

			score := rollout(state, root.me, simLimit, rand)

			backpropagate(node, score)
		}
	}

	fmt.Println("TOTAL PLAYS", root.plays)
	for m, r := range root.reward[root.me] {
		fmt.Println("CHILD", r.plays, m, r.tac)
	}

	//fmt.Println("ROOT CHILDREN NUM", len(root.children))
	//for _, child := range root.children {
	//fmt.Println("SHOULD BE DEPTH 1", child.depth)
	//for _, secondChild := range child.children {
	//fmt.Println("SHOULD BE DEPTH 2", secondChild.depth)
	//}
	//}

	fmt.Println("MAX DEPTH", maxDepth)

	//finalMove := bestMove(root, board.MeId, minPlays, t)
	finalMove := selectFinalMove(root)
	//finalMove := selectMostVisits(root)
	root.b.Print()
	fmt.Println("SELECTED MOVE", finalMove)

	return finalMove
}

func MCTSWorker(ctx context.Context, rand *rand.Rand, root *node, config MCTSConfig) rewards {
	t := determineTactic(root)
	//t := tactic(0)
	//availableWorkers := runtime.GOMAXPROCS(0)
	//availableWorkers := 1
	//fmt.Println("WORKERS", availableWorkers)

	maxDepth := 0
mctsloop:
	for {
		select {
		case <-ctx.Done():
			break mctsloop
		default:
			state := root.b.Clone()
			node := selectNode(root, state, config.MinimumSims, config.TreeDepthLimit, t)
			if node.depth > maxDepth {
				maxDepth = node.depth
			}

			score := rollout(state, root.me, config.SimLimit, rand)

			backpropagate(node, score)
		}
	}

	//fmt.Println("TOTAL PLAYS", root.plays)
	//for m, r := range root.reward[root.me] {
	//fmt.Println("CHILD", r.plays, m, r.tac)
	//}

	//fmt.Println("ROOT CHILDREN NUM", len(root.children))
	//for _, child := range root.children {
	//fmt.Println("SHOULD BE DEPTH 1", child.depth)
	//for _, secondChild := range child.children {
	//fmt.Println("SHOULD BE DEPTH 2", secondChild.depth)
	//}
	//}

	fmt.Println("MAX DEPTH", maxDepth)

	////finalMove := bestMove(root, board.MeId, minPlays, t)
	//finalMove := selectFinalMove(root)
	////finalMove := selectMostVisits(root)
	//root.b.Print()
	//fmt.Println("SELECTED MOVE", finalMove)
	rew := root.reward[root.me]
	fmt.Println("WORKER", rew)

	return root.reward[root.me]
}

func selectFinalMove(root *node) bitboard.Dir {
	t := determineTactic(root)
	me := root.me

	myRewards := root.reward[me]
	bestR := 0.0
	bestMove := moveset.Create()
	for m, r := range myRewards {
		//ensure we always have a move
		if moveset.IsEmpty(bestMove) {
			bestMove = m
		}

		score := tacticValue(r.tac, t) // + tacticValue(root.oldReward[me][m].tac, t)
		score /= r.plays               //+ root.oldReward[me][m].plays
		fmt.Println("MOVE", m, score)
		if score > bestR {
			bestR = score
			bestMove = m
		}
	}

	return bitboard.MoveSetToDir(bestMove)
}

func bestMoveByTactic(myRewards rewards, t tactic) bitboard.Dir {
	bestR := 0.0
	bestMove := moveset.Create()
	for m, r := range myRewards {
		//ensure we always have a move
		if moveset.IsEmpty(bestMove) {
			bestMove = m
		}

		score := tacticValue(r.tac, t) // + tacticValue(root.oldReward[me][m].tac, t)
		score /= r.plays               //+ root.oldReward[me][m].plays
		if score > bestR {
			bestR = score
			bestMove = m
		}
	}

	return bitboard.MoveSetToDir(bestMove)
}

func selectMostVisits(root *node) moveset.MoveSet {
	me := root.me
	myRewards := root.reward[me]
	bestR := 0.0
	bestMove := moveset.Create()

	for m, r := range myRewards {
		if r.plays > bestR {
			bestR = r.plays
			bestMove = m
		}
	}

	return bestMove
}

func determineTactic(root *node) tactic {
	me := root.me

	if root.plays == 0 {
		return survivalTac
	}

	totalSurv := 0.0
	totalFood := 0.0
	totalAgg := 0.0
	for _, r := range root.reward[me] {
		totalSurv += r.tac[survivalTac]
		totalFood += r.tac[foodTac]
		totalAgg += r.tac[aggTac]
	}
	//avgTactics := tactics{
	//totalSurv / float64(root.plays),
	//totalFood / float64(root.plays),
	//totalAgg / float64(root.plays),
	//}

	return whatDo(root.b, root.me, root.plays, totalSurv, totalFood, totalAgg)
}

func whatDo(b *bitboard.BitBoard, me string, plays, survival, food, agg float64) tactic {
	survivalThreshold := 0.6

	survivalIndicator := survival / plays
	fmt.Println("SURVIVAL INDICATOR", survivalIndicator)
	if survivalIndicator < survivalThreshold {
		fmt.Println("SURVIVAL MODE")
		return survivalTac
	}

	aliveSnakes := b.GetOpponents()
	fmt.Println("LIVING SNAKES", aliveSnakes)
	if len(b.Snakes) == 2 {
		isLargestSnake := b.GetSnake(me).Length > aliveSnakes[0].Length

		if isLargestSnake {
			fmt.Println("AGGRESSIVE MODE")
			return aggTac
		} else {
			fmt.Println("EATING MODE")
			return foodTac
		}
	}

	fmt.Println("DEFAULTING TO SURVIVAL")
	return survivalTac
}

func determineTacticFromRewards(b *bitboard.BitBoard, me string, rew rewards, totalPlays float64) tactic {
	totalSurv := 0.0
	totalFood := 0.0
	totalAgg := 0.0
	for _, r := range rew {
		totalSurv += r.tac[survivalTac]
		totalFood += r.tac[foodTac]
		totalAgg += r.tac[aggTac]
	}

	return whatDo(b, me, totalPlays, totalSurv, totalFood, totalAgg)
}

func selectNode(n *node, state *bitboard.BitBoard, minPlays, lim int, t tactic) *node {
	//fmt.Println("LOOKING AT DEPTH", n.depth)
	//state.Print()

	if state.IsGameOver() {
		//fmt.Println("GAME OVER")
		//state.Print()
		return n
	}

	if n.depth >= lim {
		//fmt.Println("HIT DEPTH LIMIT")
		//fmt.Println(n)
		return n
	}

	if len(n.children) == 0 {
		//fmt.Println("EXPANDING")
		expand(n, state)
	}

	if n.plays < float64(minPlays) {
		return n
	}

	//for _, c := range n.children {
	//fmt.Println("CHILD", c.plays, c.previousMoves, c.reward)
	//}

	selectedMoves := []bitboard.SnakeMoveSet{}
	for id, snake := range state.Snakes {
		if snake.IsAlive() {
			tactic := t
			if id != n.me {
				tactic = 0
			}
			m := bestMove(n, state, id, minPlays, tactic)
			selectedMoves = append(selectedMoves, m)
		}
	}
	//fmt.Println("SELECTED MOVES", selectedMoves)

	selectedChild := n.children[0]
	for _, child := range n.children {
		if isStateEqual(selectedMoves, child.previousMoves) {
			//fmt.Println("FOUND MATCHING CHILD")
			selectedChild = child
			break
		}
	}

	//fmt.Println(selectedChild.prevMoves)
	//state.Print()
	state.AdvanceTurn(selectedChild.prevMoves)
	//state.Print()

	//state.Print()
	//printNode(selectedChild)
	return selectNode(selectedChild, state, minPlays, lim, t)
}

//func randomUnexplored(n *node, minPlays int) *node {
//unexplored := []*node{}
//for _, c := range n.children {
//if c.plays < minPlays {
//unexplored = append(unexplored, c)
//}
//}

//return unexplored[0]
//}

func isStateEqual(a []bitboard.SnakeMoveSet, b map[string]moveset.MoveSet) bool {
	for _, m := range a {
		if m.Set != b[m.Id] {
			return false
		}
	}

	return true
}

func expand(n *node, b *bitboard.BitBoard) {
	allMoves := b.GetCartesianProductOfMoves()

	for _, m := range allMoves {
		ns := b.Clone()
		ns.AdvanceTurn(m)
		newNode := createNode(ns, m, n, n.me)
		n.children = append(n.children, newNode)
	}
}

func movesToMap(moves []bitboard.SnakeMoveSet) map[string]moveset.MoveSet {
	m := make(map[string]moveset.MoveSet, len(moves))
	for _, move := range moves {
		m[move.Id] = move.Set
	}
	return m
}

func createNode(b *bitboard.BitBoard, prevMove []bitboard.SnakeMoveSet, parent *node, me string) *node {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}

	moveMap := movesToMap(prevMove)
	rewardMatrix := createRewardMatrix(b)
	//prevRewMatrix := createRewardMatrix(ns)

	return &node{
		//b:             b,
		previousMoves: moveMap,
		prevMoves:     prevMove,
		depth:         depth,
		parent:        parent,
		reward:        rewardMatrix,
		me:            me,
		//oldReward:     prevRewMatrix,
	}
}

func createRewardMatrix(b *bitboard.BitBoard) rewardMatrix {
	rewardMatrix := make(rewardMatrix)
	for id := range b.Snakes {
		moves := moveset.Split(b.GetMoves(id).Set)
		newR := make(rewards)
		for _, m := range moves {
			moveR := reward{}
			moveR.tac = tactics{0, 0, 0}
			newR[m] = moveR
		}
		rewardMatrix[id] = newR
	}
	return rewardMatrix
}

func rollout(b *bitboard.BitBoard, meId string, turnLimit int, rand *rand.Rand) rolloutScore {
	//ns := b.Clone()
	ns := b
	me := ns.GetSnake(meId)
	curLength := 0
	if me != nil {
		curLength = me.Length
	}

	curTotalOpp := 0
	for id, snake := range ns.Snakes {
		if id != meId && snake.IsAlive() {
			curTotalOpp += 1
		}
	}

	ns.RandomPlayout(turnLimit, rand)

	me = ns.GetSnake(meId)
	foodEaten := 0
	if me != nil {
		foodEaten = me.Length - curLength
	}

	afterTotalOpp := 0
	for id, snake := range ns.Snakes {
		if id != meId && snake.IsAlive() {
			afterTotalOpp += 1
		}
	}

	oppEliminated := curTotalOpp - afterTotalOpp

	scores := make(rolloutScore)

	for id, snake := range ns.Snakes {
		survival := 0
		if snake.IsAlive() {
			survival = 1
		} else {
			survival = -1
		}

		scores[id] = tactics{
			float64(survival),
			float64(foodEaten),
			float64(oppEliminated),
		}
	}

	return scores
}

func backpropagate(n *node, rolloutReward rolloutScore) {
	n.plays += 1

	if n.parent == nil {
		return
	}

	for id, reward := range n.parent.reward {
		moveReward := reward[n.previousMoves[id]]
		if r, ok := rolloutReward[id]; ok {
			moveReward.tac[0] += r[0]
			moveReward.tac[1] += r[1]
			moveReward.tac[2] += r[2]
		}
		moveReward.plays += 1
		n.parent.reward[id][n.previousMoves[id]] = moveReward
	}

	//for id, r := range rewards {
	//parentR := n.parent.reward[id]
	//parentT := parentR[n.previousMoves[id]]

	//parentT.tac[0] += r[0]
	//parentT.tac[1] += r[1]
	//parentT.tac[2] += r[2]
	////parentT.tac[0] = math.Max(parentT.tac[0], r[0])
	////parentT.tac[1] = math.Max(parentT.tac[1], r[1])
	////parentT.tac[2] = math.Max(parentT.tac[2], r[2])

	//parentT.plays += 1
	//n.parent.reward[id][n.previousMoves[id]] = parentT
	//}

	backpropagate(n.parent, rolloutReward)
}

func bestMove(n *node, b *bitboard.BitBoard, id string, minPlays int, t tactic) bitboard.SnakeMoveSet {
	bestMove := moveset.Create()
	bestUtc := 0.0

	for dir := range n.reward[id] {
		v := uct(n, minPlays, t, id, dir)
		if v > bestUtc {
			bestUtc = v
			bestMove = dir
		}
	}
	//fmt.Println("BEST UTC", bestMove, bestUtc)

	//fmt.Println("SELECTED MOVE", children[0].previousMoves[id], uct(children[0], minPlays, t, id))
	return bitboard.SnakeMoveSet{Id: id, Set: bestMove}
}

func uct(n *node, minPlays int, t tactic, id string, move moveset.MoveSet) float64 {
	if n.reward[id][move].plays < float64(minPlays) {
		return math.MaxFloat64
	}
	score := tacticValue(n.reward[id][move].tac, t)
	//decayedScore := tacticValue(n.oldReward[id][move].tac, t)
	//score += decayedScore

	plays := n.reward[id][move].plays
	//decayedPlays := n.oldReward[id][move].plays
	//plays += decayedPlays

	u := calculateUCT(n, score, plays, minPlays)
	return u
}

func calculateUCT(n *node, score float64, plays float64, minPlays int) float64 {
	//if n.plays < minPlays {
	//return math.MaxFloat64
	//}

	parentPlays := n.plays
	explorConst := math.Sqrt(2)

	p := float64(plays)
	exploitation := score / p
	exploration := explorConst * math.Sqrt(math.Log(parentPlays)/p)

	return exploitation + exploration
}

//func uctTuned(n *node, minPlays int, t tactic, id board.SnakeId, move board.Move) float64 {
//score := tacticValue(n.reward[id][move].tac, 0)
//plays := n.reward[id][move].plays
////explorationConstant := math.Sqrt(2)
////alpha := float64(0.1)

//numParentSims := float64(n.plays)
////heuristic := payoff.heuristic[move]
//variance := math.Pow(score, 2) / plays
//mean := nodeMean(n, id, t)
//v := (0.5 * variance) - mean + math.Sqrt((2*math.Log(numParentSims))/plays)

////exploitation := (1-alpha)*(score/plays) + alpha*heuristic
//exploitation := score / plays
//exploration := math.Sqrt((math.Log(numParentSims) / plays) * math.Min(1/4, v))

////fmt.Println("ucb ", exploitation+exploration)
//return exploitation + exploration
//}

//func nodeMean(node *node, id board.SnakeId, t tactic) float64 {
//nReward := node.reward[id]
//totalScore := 0.0
//totalSquared := 0.0
//for _, s := range nReward {
//tac := tacticValue(s.tac, t)
//totalScore += tac
//totalSquared += math.Pow(float64(tac), 2)
//}
//mean := float64(totalScore / node.plays)
//return mean
//}

func tacticValue(tactics tactics, t tactic) float64 {
	if t == 0 {
		return tactics[t]
	}

	adjTac := tactics[t] * tactics[0]
	if adjTac == 0 {
		return tactics[0]
	}

	return adjTac
}

//func randomRollout(fb *board.FastBoard, limit int) rolloutScore {
////fmt.Println("START ROLLOUT")
////fb.Print()

//moves := make(map[board.SnakeId]board.Move, len(fb.Lengths))
//turn := 0
//foodBefore := len(fb.GetAllFood())
//foodEaten := 0

//liveSnakes := []board.SnakeId{}
//for id := range fb.Heads {
//if fb.IsSnakeAlive(id) && id != board.MeId {
//liveSnakes = append(liveSnakes, id)
//}
//}

//for !fb.IsGameOver() && fb.IsSnakeAlive(board.MeId) && turn < limit {
////fb.Print()
////fmt.Println("TURN", turn)
//for id := range fb.Lengths {
//if !fb.IsSnakeAlive(id) {
//moves[id] = ""
//continue
//}
//move := randomAgentMove(fb, id)
//moves[id] = move.Dir
//}
//fb.AdvanceBoard(moves)
//turn += 1

//hasEaten := fb.Healths[board.MeId] == 100 && fb.Lengths[board.MeId] > 3

//if hasEaten {
//foodEaten += 1
//}
//}
////fb.Print()

//totalDeadSnakes := 0
//for _, id := range liveSnakes {
//if !fb.IsSnakeAlive(id) {
//totalDeadSnakes += 1
//}
//}

//snakesKilledScore := float64(totalDeadSnakes) / float64(len(liveSnakes)+1)
//foodScore := float64(foodEaten) / float64(foodBefore+1)

//score := make(rolloutScore)
//for id := range fb.Heads {
//if fb.IsSnakeAlive(id) {
//score[id] = tactics{1, 0, 0}
//} else {
//score[id] = tactics{0, 0, 0}
//}
//}

//meScore := score[board.MeId]
//meScore[1] = foodScore
//meScore[2] = snakesKilledScore
//score[board.MeId] = meScore

////fmt.Println("END ROLLOUT", rewards)
////return tactics{survivalScore, foodScore, snakesKilledScore}
//return score
