package brain

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"

	g "github.com/nosnaws/tiam/game"
)

type Node struct {
	board         *g.FastBoard
	children      []*Node
	parent        *Node
	plays         int
	moveSet       map[g.SnakeId]g.SnakeMove
	possibleMoves map[g.SnakeId][]g.SnakeMove
	payoffs       map[g.SnakeId]Payoff
	heuristics    Heuristics
	depth         int
}

type Payoff struct {
	plays  map[g.Move]int
	scores map[g.Move]int
	//heuristic map[g.Move]float64
}

type SnakeScore struct {
	id    g.SnakeId
	value int
	move  g.Move
	//heuristic float64
}

type MCTSConfig struct {
	ExplorationConstant float64
	AlphaConstant       float64
	VoronoiWeighting    float64
	FoodWeightA         float64
	FoodWeightB         float64
	BigSnakeReward      float64
}

type MASTTable struct {
	total    int
	moveWins map[g.Move]int
}

type cacheEntry struct {
	plays   int
	payoffs map[g.SnakeId]Payoff
}

type MctsGame struct {
	tt    TranspositionTable
	cache map[BoardHash]cacheEntry
}

func CreateMctsGame(h, w int) MctsGame {
	return MctsGame{
		tt:    InitializeTranspositionTable(h, w),
		cache: make(map[BoardHash]cacheEntry),
	}
}

func (mc *MctsGame) MCTS(board *g.FastBoard, config *MCTSConfig, txn *newrelic.Transaction) g.SnakeMove {
	segment := txn.StartSegment("MCTS")
	defer segment.End()

	fakeMoveSet := make(map[g.SnakeId]g.SnakeMove)
	root := createNode(fakeMoveSet, board)
	//root.children = mc.createChildren(root)

	duration, err := time.ParseDuration("350ms")
	if err != nil {
		panic("could not parse duration")
	}

	// create MAST table
	//mast := initMAST(*root.board)

	maxDepth := 0
loop:
	for timeout := time.After(duration); ; {
		select {
		case <-timeout:
			break loop
		default:
			depth := mctsIteration(root)

			if maxDepth < depth {
				maxDepth = depth
			}

		}
	}

	bestMove := selectFinalMove(root)
	//bestMove := bestMoveUTC(root, fastGame.MeId)
	printNode(root)
	//for _, child := range root.children {
	//printNode()
	//for _, c := range child.children {
	//printNode(c)
	//}
	//}
	log.Println("# Selected #")
	log.Println(bestMove)
	log.Println("Total plays: ", root.plays)
	log.Println("Max depth: ", maxDepth)
	addAttributes(txn, root, bestMove, maxDepth)
	return bestMove
}

func mctsIteration(root *Node) int {
	// initalize the tree
	if root.plays == 0 {
		expandNode(root)
		root.heuristics = getHeuristics(root)
	}

	node := selectNode(root)
	fmt.Println("selecting node", node.moveSet)
	fmt.Println("selecting depth", node.depth)
	node.board.Print()

	var score Rewards
	child := getRandomUnexploredChild(node)
	fmt.Println("random child to playout ", child.moveSet)
	if !child.board.IsGameOver() {
		fmt.Println("game is not over, expanding")

		expandNode(child)
		child.heuristics = getHeuristics(child)

		score = simulateNode(child)
		fmt.Println("playout result", score)
		//node.plays += 1
		//updateMAST(mast, score)

	} else {
		score = getRewards(child, child.board)
		fmt.Println("game is over, getting rewards", score)
	}

	backpropagate(node, score)

	return node.depth
}

type Rewards map[g.SnakeId]SnakeScore
type Heuristics map[g.SnakeId]float64

func getRewards(node *Node, playedOutBoard *g.FastBoard) Rewards {
	rewards := make(Rewards)

	for id := range playedOutBoard.Heads {
		val := 0

		if playedOutBoard.IsSnakeAlive(id) {
			val = 1
		}

		rewards[id] = SnakeScore{
			id:    id,
			value: val,
			move:  node.moveSet[id].Dir,
		}
	}

	return rewards
}

func getHeuristics(node *Node) Heuristics {
	maxH := make(Heuristics)
	for id := range node.board.Heads {
		maxH[id] = 0
	}

	for _, child := range node.children {
		bh := calculateBoardHeuristic(*child.board)

		for id, h := range maxH {
			h = math.Max(h, bh[id])
			maxH[id] = h
		}
	}

	return maxH
}

func selectNode(node *Node) *Node {
	fmt.Println("looking node depth", node.depth)
	fmt.Println("looking node rewards", node.payoffs)
	fmt.Println("looking node heuristic", node.heuristics)
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node))
}

func isLeafNode(node *Node) bool {
	if node.plays < len(node.children) {
		return true
	}

	return false
}

func expandNode(node *Node) {
	node.children = createChildren(node)
}

func simulateNode(node *Node) Rewards {
	ns := node.board.Clone()
	ns.RandomRollout()
	return getRewards(node, &ns)

	//if ns.IsSnakeAlive(fastGame.MeId) {
	//ns.RandomRollout()
	//turn := 0
	//moves := make(map[g.SnakeId]g.Move)
	//for !ns.IsGameOver() && turn < 10 {

	//for id := range ns.Lengths {
	//if !ns.IsSnakeAlive(id) {
	//moves[id] = ""
	//continue
	//}

	////sMoves := b.GetMovesForSnake(id)
	////randomMove := sMoves[rand.Intn(len(sMoves))]
	////moves = append(moves, randomMove)
	//moves[id] = selectSimMove(ns, id, mast[id])
	//}
	////turn += 1
	//ns.AdvanceBoard(moves)
	//turn += 1
	//}

	//heuristic := calculateBoardHeuristic(*node.board)
	//}

	//nodeHeuristic := calculateNodeHeuristic(node, g.MeId, config)

	//isDraw := true
	//for id := range ns.Lengths {
	//if ns.IsSnakeAlive(id) {
	//isDraw = false
	//}
	//}

	//scores := make(map[g.SnakeId]SnakeScore, len(ns.Heads))
	//for id := range ns.Lengths {
	////snakeHeuristic := nodeHeuristic
	////if id != g.MeId {
	////snakeHeuristic = -snakeHeuristic
	////}

	//score := SnakeScore{
	//id:        id,
	//value:     0,
	//move:      node.moveSet[id].Dir,
	//heuristic: heuristic[id],
	//}

	//if ns.IsSnakeAlive(id) {
	//score.value = 1
	//}

	////if isDraw {
	////score.value = 0
	////} else if ns.IsSnakeAlive(id) {
	////score.value = 1
	////} else {
	////score.value = -1
	////}
	//scores[id] = score
	//}

	//return scores
}

func backpropagate(node *Node, rewards Rewards) {
	if node == nil {
		return
	}

	pastMovesWithScore := make(map[g.SnakeId]SnakeScore)
	for id := range node.board.Lengths {
		if payoff, ok := node.payoffs[id]; ok {

			reward := rewards[id]

			payoff.plays[reward.move] += 1
			payoff.scores[reward.move] += reward.value

			//h := payoff.heuristic[reward.move]
			//payoff.heuristic[reward.move] = math.Max(reward.heuristic, h)

			node.payoffs[id] = payoff
		}

		if _, ok := node.moveSet[id]; ok {
			val := 0
			if _, ok := rewards[id]; ok {
				val = rewards[id].value
			}
			pastMovesWithScore[id] = SnakeScore{
				id:    id,
				value: val,
				move:  node.moveSet[id].Dir,
			}
		}
	}

	node.plays += 1

	for _, child := range node.children {
		for id, h := range child.heuristics {
			node.heuristics[id] = math.Max(h, node.heuristics[id])
		}
	}

	//boardHash := HashBoard(mc.tt, *node.board)
	//mc.cache[boardHash] = cacheEntry{
	//plays:   node.plays,
	//payoffs: node.payoffs,
	//}

	backpropagate(node.parent, pastMovesWithScore)
}

func getRandomUnexploredChild(node *Node) *Node {
	var unexplored []*Node
	for _, child := range node.children {
		if child.plays == 0 {
			unexplored = append(unexplored, child)
			//return child
		}
	}

	//return nil
	if len(unexplored) > 0 {
		return ShuffleNodes(unexplored)[0]
	}
	return nil
}

func createChildren(node *Node) []*Node {
	productOfMoves := g.GetCartesianProductOfMoves(*node.board)

	var children []*Node
	for _, moveSet := range productOfMoves {
		cs := node.board.Clone()
		cs.AdvanceBoard(movesToMap(moveSet))

		moves := make(map[g.SnakeId]g.SnakeMove)
		for _, m := range moveSet {
			moves[m.Id] = m
		}

		childNode := createNode(moves, &cs)
		childNode.parent = node
		childNode.depth = childNode.parent.depth + 1

		children = append(children, childNode)
	}

	return children
}

func movesToMap(moves []g.SnakeMove) map[g.SnakeId]g.Move {
	m := make(map[g.SnakeId]g.Move, len(moves))
	for _, move := range moves {
		m[move.Id] = move.Dir
	}
	return m
}

func createNode(moveSet map[g.SnakeId]g.SnakeMove, board *g.FastBoard) *Node {
	possibleMoves := make(map[g.SnakeId][]g.SnakeMove)
	payoffs := make(map[g.SnakeId]Payoff)
	for id := range board.Lengths {
		if !board.IsSnakeAlive(id) {
			continue
		}
		moves := board.GetMovesForSnake(id)
		possibleMoves[id] = moves

		payoffs[id] = createPayoff(moves)
	}

	node := Node{
		possibleMoves: possibleMoves,
		board:         board,
		payoffs:       payoffs,
		moveSet:       moveSet,
	}

	//nodeHash := HashBoard(mc.tt, *board)
	//if cached, ok := mc.cache[nodeHash]; ok {
	//node.plays = cached.plays
	//node.payoffs = cached.payoffs
	//}

	return &node
}

func createPayoff(moves []g.SnakeMove) Payoff {
	plays := make(map[g.Move]int, len(moves))
	scores := make(map[g.Move]int, len(moves))
	//heuristic := make(map[g.Move]float64, len(moves))

	for _, m := range moves {
		plays[m.Dir] = 0
		scores[m.Dir] = 0
		//heuristic[m.Dir] = 0
	}
	return Payoff{plays: plays, scores: scores}
}

func ShuffleNodes(nodes []*Node) []*Node {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]*Node, len(nodes))
	perm := r.Perm(len(nodes))
	for i, randIndex := range perm {
		ret[i] = nodes[randIndex]
	}
	return ret
}

func ShuffleMoves(moves []g.SnakeMove) []g.SnakeMove {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]g.SnakeMove, len(moves))
	perm := r.Perm(len(moves))
	for i, randIndex := range perm {
		ret[i] = moves[randIndex]
	}
	return ret
}

func selectFinalMove(node *Node) g.SnakeMove {
	moves := node.possibleMoves[g.MeId]
	sort.Slice(moves, func(a, b int) bool {
		return moveSecureness(node, g.MeId, moves[a]) > moveSecureness(node, g.MeId, moves[b])
	})

	return moves[0]
}

func moveSecureness(node *Node, player g.SnakeId, move g.SnakeMove) float64 {
	//score := float64(node.payoffs[player].scores[move.Dir])
	plays := float64(node.payoffs[player].plays[move.Dir])

	return plays
}

func bestUTC(node *Node) *Node {
	var moveSet []g.SnakeMove
	for id := range node.board.Lengths {
		if node.board.IsSnakeAlive(id) {
			bestMove := bestMoveUTC(node, id)
			moveSet = append(moveSet, bestMove)
		}
	}

	for _, child := range node.children {
		if isStateEqual(moveSet, child.moveSet) {
			return child
		}
	}

	return nil
}

func isStateEqual(a []g.SnakeMove, b map[g.SnakeId]g.SnakeMove) bool {
	equal := true
	for _, m := range a {
		if m.Dir != b[m.Id].Dir {
			equal = false
		}
	}

	return equal
}

func bestMoveUTC(node *Node, id g.SnakeId) g.SnakeMove {
	moves := node.possibleMoves[id]
	fmt.Println("getting best move for ", id)
	sort.Slice(moves, func(a, b int) bool {
		ucbA := calculateUCB(node, id, moves[a].Dir)
		ucbB := calculateUCB(node, id, moves[b].Dir)
		fmt.Println("comparing a", moves[a].Dir, ucbA)
		fmt.Println("comparing b", moves[b].Dir, ucbB)
		return ucbA > ucbB
	})

	fmt.Println("UTC ", id, moves)

	return moves[0]
}

func calculateUCB(node *Node, id g.SnakeId, move g.Move) float64 {
	payoff := node.payoffs[id]
	//explorationConstant := math.Sqrt(config.ExplorationConstant)
	explorationConstant := math.Sqrt(2)

	//alpha := float64(config.AlphaConstant)
	alpha := 0.4

	numParentSims := float64(node.plays)
	score := float64(payoff.scores[move])
	plays := float64(payoff.plays[move])
	heuristic := node.heuristics[id]

	exploitation := (1-alpha)*(score/plays) + alpha*heuristic
	//exploitation := score / plays
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/plays)

	return exploitation + exploration
}

func calculateNodeHeuristic(node *Node, id g.SnakeId, config *MCTSConfig) float64 {
	health := float64(node.board.Healths[id])

	if !node.board.IsSnakeAlive(id) {
		return -1.0
	}

	//isLargestSnake := true
	//for sId, l := range node.board.Lengths {
	//if sId != id && l >= node.board.Lengths[id] {
	//isLargestSnake = false
	//}
	//}

	voronoi := g.Voronoi(node.board, id)

	total := 0.0
	//if isLargestSnake {
	//total += config.BigSnakeReward
	//}

	total += config.FoodWeightA * ((health - float64(voronoi.FoodDepth[id])) / config.FoodWeightB)
	total += config.VoronoiWeighting * float64(voronoi.Score[id])

	return total / math.Sqrt(100+math.Pow(total, 2))
}

func calculateBoardHeuristic(b g.FastBoard) map[g.SnakeId]float64 {
	scores := make(map[g.SnakeId]float64, len(b.Heads))

	//scores[largestSnake] += 10

	voronoi := g.Voronoi(&b, g.MeId)

	for id, score := range voronoi.Score {
		scores[id] += 1 * float64(score) / 2
	}

	for id, f := range voronoi.FoodDepth {
		scores[id] += -float64(f)
	}

	for id, s := range scores {
		scores[id] = s / math.Sqrt(100+math.Pow(s, 2))
	}

	for id := range b.Lengths {
		//if id != g.MeId && !b.IsSnakeAlive(g.MeId) {
		//scores[id] = math.MaxFloat64
		//continue
		//}
		if !b.IsSnakeAlive(id) {
			scores[id] = -math.MaxFloat64
		} else if b.IsGameOver() {
			scores[id] = math.MaxFloat64
		}
	}

	return scores
}

func updateMAST(mast map[g.SnakeId]MASTTable, score map[g.SnakeId]SnakeScore) {
	for id, s := range score {
		sMast := mast[id]

		sMast.moveWins[s.move] += s.value
		sMast.total += 1
		mast[id] = sMast
	}
}

func initMAST(b g.FastBoard) map[g.SnakeId]MASTTable {
	m := make(map[g.SnakeId]MASTTable)

	for id := range b.Heads {
		mast := m[id]
		mast.moveWins = make(map[g.Move]int)
		mast.moveWins[g.Left] = 0
		mast.moveWins[g.Up] = 0
		mast.moveWins[g.Down] = 0
		mast.moveWins[g.Right] = 0
		mast.total = 1
		m[id] = mast
	}

	return m
}

func addAttributes(txn *newrelic.Transaction, root *Node, selected g.SnakeMove, maxDepth int) {
	if txn != nil {
		txn.AddAttribute("totalPlays", root.plays)
		txn.AddAttribute("selectedMove", selected.Dir)
		txn.AddAttribute("maxDepth", maxDepth)
	}
}

func printNode(node *Node) {
	log.Println("#############")
	node.board.Print()
	log.Println("Depth", node.depth)
	log.Println("Total plays", node.plays)
	log.Println("Heuristics", node.heuristics)

	for id, payoff := range node.payoffs {
		log.Println("Player", id)
		log.Println("Health", node.board.Healths[id])
		log.Println("Length", node.board.Lengths[id])
		log.Println("Plays", payoff.plays)
		log.Println("Scores", payoff.scores)
	}

	for _, child := range node.children {
		printChild(child)
	}
}

func printChild(node *Node) {
	log.Printf("-- depth:%d moves: %v --", node.depth, node.moveSet)
	log.Println("Total plays", node.plays)
	log.Println("Heuristics", node.heuristics)
	for id, payoff := range node.payoffs {
		log.Println("Player", id)
		log.Println("Plays", payoff.plays)
		log.Println("Scores", payoff.scores)
		//log.Println("Heuristics", payoff.heuristic)
	}
}

func selectSimMove(b g.FastBoard, id g.SnakeId, mast MASTTable) g.Move {
	totalAvg := 1
	for _, w := range mast.moveWins {
		totalAvg += w / mast.total
	}

	moves := ShuffleMoves(b.GetMovesForSnake(id))

	r := rand.Float32()
	if r > 0.6 {
		return moves[0].Dir
	}

	sort.Slice(moves, func(i, j int) bool {
		return mast.moveWins[moves[i].Dir]/totalAvg > mast.moveWins[moves[j].Dir]/totalAvg
	})

	return moves[0].Dir
}
