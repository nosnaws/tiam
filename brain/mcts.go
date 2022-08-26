package brain

import (
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"

	fastGame "github.com/nosnaws/tiam/game"
)

type Node struct {
	board         *fastGame.FastBoard
	children      []*Node
	parent        *Node
	plays         int
	moveSet       map[fastGame.SnakeId]fastGame.SnakeMove
	possibleMoves map[fastGame.SnakeId][]fastGame.SnakeMove
	payoffs       map[fastGame.SnakeId]Payoff
	depth         int
}

type Payoff struct {
	plays     map[fastGame.Move]int
	scores    map[fastGame.Move]int
	heuristic map[fastGame.Move]float64
}

type SnakeScore struct {
	id        fastGame.SnakeId
	value     int
	move      fastGame.Move
	heuristic float64
}

type MCTSConfig struct {
	ExplorationConstant float64
	AlphaConstant       float64
	VoronoiWeighting    float64
	FoodWeightA         float64
	FoodWeightB         float64
	BigSnakeReward      float64
}

func addAttributes(txn *newrelic.Transaction, root *Node, selected fastGame.SnakeMove, maxDepth int) {
	if txn != nil {
		txn.AddAttribute("totalPlays", root.plays)
		txn.AddAttribute("selectedMove", selected.Dir)
		txn.AddAttribute("maxDepth", maxDepth)
	}
}

func MCTS(board *fastGame.FastBoard, config *MCTSConfig, txn *newrelic.Transaction) fastGame.SnakeMove {
	segment := txn.StartSegment("MCTS")
	defer segment.End()

	fakeMoveSet := make(map[fastGame.SnakeId]fastGame.SnakeMove)
	root := createNode(fakeMoveSet, board)
	root.children = createChildren(root)

	duration, err := time.ParseDuration("350ms")
	if err != nil {
		panic("could not parse duration")
	}

	maxDepth := 0
loop:
	for timeout := time.After(duration); ; {
		select {
		case <-timeout:
			break loop
		default:
			node := selectNode(root, config)

			child := expandNode(node)

			score := simulateNode(child, config)
			child.plays += 1
			if maxDepth < child.depth {
				maxDepth = child.depth
			}

			backpropagate(node, score)
		}
	}

	bestMove := selectFinalMove(root)
	//bestMove := bestMoveUTC(root, fastGame.MeId)
	printNode(root)
	//for _, child := range root.children {
	//printNode(child)
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

func selectNode(node *Node, config *MCTSConfig) *Node {
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node, config), config)
}

func printNode(node *Node) {
	log.Println("#############")
	node.board.Print()
	log.Println("Depth", node.depth)
	log.Println("Total plays", node.plays)

	for id, payoff := range node.payoffs {
		log.Println("Player", id)
		log.Println("Health", node.board.Healths[id])
		log.Println("Length", node.board.Lengths[id])
		log.Println("Plays", payoff.plays)
		log.Println("Scores", payoff.scores)
		log.Println("Heuristics", payoff.heuristic)
	}

	for _, child := range node.children {
		printChild(child)
	}
}

func printChild(node *Node) {
	log.Printf("-- depth:%d moves: %v --", node.depth, node.moveSet)
	log.Println("Total plays", node.plays)
	for id, payoff := range node.payoffs {
		log.Println("Player", id)
		log.Println("Plays", payoff.plays)
		log.Println("Scores", payoff.scores)
		log.Println("Heuristics", payoff.heuristic)
	}
}

func expandNode(node *Node) *Node {
	if node.board.IsGameOver() {
		return node
	}

	if len(node.children) == 0 {
		node.children = createChildren(node)
	}

	return getRandomUnexploredChild(node)
}

func simulateNode(node *Node, config *MCTSConfig) map[fastGame.SnakeId]SnakeScore {

	ns := node.board.Clone()

	//if ns.IsSnakeAlive(fastGame.MeId) {
	ns.RandomRollout()
	//}

	nodeHeuristic := calculateNodeHeuristic(node, fastGame.MeId, config)

	//isDraw := true
	//for id := range ns.Lengths {
	//if ns.IsSnakeAlive(id) {
	//isDraw = false
	//}
	//}

	scores := make(map[fastGame.SnakeId]SnakeScore, len(ns.Heads))
	for id := range ns.Lengths {
		snakeHeuristic := nodeHeuristic
		if id != fastGame.MeId {
			snakeHeuristic = -snakeHeuristic
		}

		score := SnakeScore{
			id:        id,
			value:     0,
			move:      node.moveSet[id].Dir,
			heuristic: snakeHeuristic,
		}

		if ns.IsSnakeAlive(id) {
			score.value = 1
		}

		//if isDraw {
		//score.value = 0
		//} else if ns.IsSnakeAlive(id) {
		//score.value = 1
		//} else {
		//score.value = -1
		//}
		scores[id] = score
	}

	return scores
}

func backpropagate(node *Node, scores map[fastGame.SnakeId]SnakeScore) {
	if node == nil {
		return
	}

	pastMovesWithScore := make(map[fastGame.SnakeId]SnakeScore)
	node.plays += 1
	for id := range node.board.Lengths {
		if payoff, ok := node.payoffs[id]; ok {
			score := scores[id]
			payoff.plays[score.move] += 1

			payoff.scores[score.move] += score.value

			h := payoff.heuristic[score.move]
			payoff.heuristic[score.move] = math.Max(score.heuristic, h)

			node.payoffs[id] = payoff
		}

		if _, ok := node.moveSet[id]; ok {
			val := 0
			if _, ok := scores[id]; ok {
				val = scores[id].value
			}
			pastMovesWithScore[id] = SnakeScore{
				id:        id,
				value:     val,
				move:      node.moveSet[id].Dir,
				heuristic: scores[id].heuristic,
			}
		}
	}

	backpropagate(node.parent, pastMovesWithScore)
}

func isLeafNode(node *Node) bool {
	if len(node.children) == 0 {
		return true
	}

	for _, n := range node.children {
		if n.plays < 1 {
			return true
		}
	}

	return false
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
	return Shuffle(unexplored)[0]
}

func createChildren(node *Node) []*Node {
	productOfMoves := fastGame.GetCartesianProductOfMoves(*node.board)

	var children []*Node
	for _, moveSet := range productOfMoves {
		cs := node.board.Clone()
		cs.AdvanceBoard(movesToMap(moveSet))

		moves := make(map[fastGame.SnakeId]fastGame.SnakeMove)
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

func movesToMap(moves []fastGame.SnakeMove) map[fastGame.SnakeId]fastGame.Move {
	m := make(map[fastGame.SnakeId]fastGame.Move, len(moves))
	for _, move := range moves {
		m[move.Id] = move.Dir
	}
	return m
}

func createNode(moveSet map[fastGame.SnakeId]fastGame.SnakeMove, board *fastGame.FastBoard) *Node {
	possibleMoves := make(map[fastGame.SnakeId][]fastGame.SnakeMove)
	payoffs := make(map[fastGame.SnakeId]Payoff)
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

	return &node
}

func createPayoff(moves []fastGame.SnakeMove) Payoff {
	plays := make(map[fastGame.Move]int, len(moves))
	scores := make(map[fastGame.Move]int, len(moves))
	heuristic := make(map[fastGame.Move]float64, len(moves))

	for _, m := range moves {
		plays[m.Dir] = 0
		scores[m.Dir] = 0
		heuristic[m.Dir] = 0
	}
	return Payoff{plays: plays, scores: scores, heuristic: heuristic}
}

func Shuffle(nodes []*Node) []*Node {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]*Node, len(nodes))
	perm := r.Perm(len(nodes))
	for i, randIndex := range perm {
		ret[i] = nodes[randIndex]
	}
	return ret
}

func selectFinalMove(node *Node) fastGame.SnakeMove {
	moves := node.possibleMoves[fastGame.MeId]
	sort.Slice(moves, func(a, b int) bool {
		return moveSecureness(node, fastGame.MeId, moves[a]) > moveSecureness(node, fastGame.MeId, moves[b])
	})

	return moves[0]
}

func moveSecureness(node *Node, player fastGame.SnakeId, move fastGame.SnakeMove) float64 {
	//score := float64(node.payoffs[player].scores[move.Dir])
	plays := float64(node.payoffs[player].plays[move.Dir])

	return plays
}

func bestUTC(node *Node, config *MCTSConfig) *Node {
	var moveSet []fastGame.SnakeMove
	for id := range node.board.Lengths {
		if node.board.IsSnakeAlive(id) {
			bestMove := bestMoveUTC(node, id, config)
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

func isStateEqual(a []fastGame.SnakeMove, b map[fastGame.SnakeId]fastGame.SnakeMove) bool {
	equal := true
	for _, m := range a {
		if m.Dir != b[m.Id].Dir {
			equal = false
		}
	}

	return equal
}

func bestMoveUTC(node *Node, id fastGame.SnakeId, config *MCTSConfig) fastGame.SnakeMove {
	moves := node.possibleMoves[id]
	sort.Slice(moves, func(a, b int) bool {
		return calculateUCB(node, id, moves[a].Dir, config) > calculateUCB(node, id, moves[b].Dir, config)
	})

	return moves[0]
}

func calculateUCB(node *Node, id fastGame.SnakeId, move fastGame.Move, config *MCTSConfig) float64 {
	payoff := node.payoffs[id]
	explorationConstant := math.Sqrt(config.ExplorationConstant)
	alpha := float64(config.AlphaConstant)

	numParentSims := float64(node.plays)
	score := float64(payoff.scores[move])
	plays := float64(payoff.plays[move])
	heuristic := payoff.heuristic[move]

	exploitation := (1-alpha)*(score/plays) + alpha*heuristic
	//exploitation := score / plays
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/plays)

	return exploitation + exploration
}

func calculateNodeHeuristic(node *Node, id fastGame.SnakeId, config *MCTSConfig) float64 {
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

	voronoi := fastGame.Voronoi(node.board, id)

	total := 0.0
	//if isLargestSnake {
	//total += config.BigSnakeReward
	//}

	total += config.FoodWeightA * ((health - float64(voronoi.FoodDepth)) / config.FoodWeightB)
	total += config.VoronoiWeighting * float64(voronoi.Score)

	return total / math.Sqrt(100+math.Pow(total, 2))
}
