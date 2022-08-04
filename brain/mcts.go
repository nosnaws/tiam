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

func addAttributes(txn *newrelic.Transaction, root *Node, selected fastGame.SnakeMove) {
	if txn != nil {
		txn.AddAttribute("totalPlays", root.plays)
		txn.AddAttribute("selectedMove", selected.Dir)
	}
}

func MCTS(board *fastGame.FastBoard, txn *newrelic.Transaction) fastGame.SnakeMove {
	fakeMoveSet := make(map[fastGame.SnakeId]fastGame.SnakeMove)
	root := createNode(fakeMoveSet, board)
	root.children = createChildren(root)

	duration, err := time.ParseDuration("350ms")
	if err != nil {
		panic("could not parse duration")
	}

loop:
	for timeout := time.After(duration); ; {
		select {
		case <-timeout:
			break loop
		default:
			node := selectNode(root)
			child := expandNode(node)

			score := simulateNode(child)
			child.plays += 1

			backpropagate(node, score)
		}
	}

	bestMove := selectFinalMove(root)
	log.Println("# Selected #")
	log.Println(bestMove)
	log.Println("Total plays: ", root.plays)
	addAttributes(txn, root, bestMove)
	return bestMove
}

func selectNode(node *Node) *Node {
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node))
}

func printNode(node *Node) {
	log.Println("#############")
	log.Println("Depth", node.depth)
	log.Println("payoffs", node.payoffs)
	log.Println("#############")
}

func expandNode(node *Node) *Node {
	if node.board.IsGameOver() {
		return node
	}

	if len(node.children) == 0 {
		children := createChildren(node)
		if len(children) == 0 {
			node.board.Print()
		}

		node.children = children
	}

	return getRandomUnexploredChild(node)
}

func simulateNode(node *Node) map[fastGame.SnakeId]SnakeScore {

	ns := node.board.Clone()
	ns.RandomRollout()

	nodeHeuristic := calculateNodeHeuristic(node, fastGame.MeId)

	scores := make(map[fastGame.SnakeId]SnakeScore, len(ns.Heads))
	for id, l := range ns.Lengths {
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

		if l > 0 {
			score.value = 1
		}
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

		if _, ok := node.payoffs[id]; ok {
			score := scores[id]
			node.payoffs[id].plays[score.move] += 1

			node.payoffs[id].scores[score.move] += score.value

			h := node.payoffs[id].heuristic[score.move]
			node.payoffs[id].heuristic[score.move] = math.Max(h, score.heuristic)
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
		}
	}

	return Shuffle(unexplored)[0]
}

func createChildren(node *Node) []*Node {
	productOfMoves := fastGame.GetCartesianProductOfMoves(*node.board)

	var children []*Node
	for _, moveSet := range productOfMoves {
		cs := node.board.Clone()
		cs.AdvanceBoard(moveSet)

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

func createNode(moveSet map[fastGame.SnakeId]fastGame.SnakeMove, board *fastGame.FastBoard) *Node {
	possibleMoves := make(map[fastGame.SnakeId][]fastGame.SnakeMove)
	payoffs := make(map[fastGame.SnakeId]Payoff)
	for id, length := range board.Lengths {
		isAlive := length > 0
		if !isAlive {
			continue
		}
		moves := board.GetMovesForSnake(id)
		possibleMoves[id] = moves
		payoffs[id] = createPayoff(moves)
	}

	return &Node{
		possibleMoves: possibleMoves,
		board:         board,
		payoffs:       payoffs,
		moveSet:       moveSet,
	}
}

func createPayoff(moves []fastGame.SnakeMove) Payoff {
	plays := make(map[fastGame.Move]int)
	scores := make(map[fastGame.Move]int)
	heuristics := make(map[fastGame.Move]float64)

	for _, m := range moves {
		plays[m.Dir] = 0
		scores[m.Dir] = 0
		heuristics[m.Dir] = 0
	}
	return Payoff{plays: plays, scores: scores, heuristic: heuristics}
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
	payoff := node.payoffs[fastGame.MeId]
	sort.Slice(moves, func(a, b int) bool {
		return payoff.scores[moves[a].Dir] > payoff.scores[moves[b].Dir]
	})

	return moves[0]
}

func bestUTC(node *Node) *Node {
	var moveSet []fastGame.SnakeMove
	for id, h := range node.board.Healths {
		if h > 0 {
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

func isStateEqual(a []fastGame.SnakeMove, b map[fastGame.SnakeId]fastGame.SnakeMove) bool {
	equal := true
	for _, m := range a {
		if m.Dir != b[m.Id].Dir {
			equal = false
		}
	}

	return equal
}

func bestMoveUTC(node *Node, id fastGame.SnakeId) fastGame.SnakeMove {
	moves := node.possibleMoves[id]
	sort.Slice(moves, func(a, b int) bool {
		return calculateUCB(node, id, moves[a].Dir) > calculateUCB(node, id, moves[b].Dir)
	})

	return moves[0]
}

func calculateUCB(node *Node, id fastGame.SnakeId, move fastGame.Move) float64 {
	payoff := node.payoffs[id]
	explorationConstant := math.Sqrt(2)
	alpha := float64(0.1)

	numParentSims := float64(node.plays)
	score := payoff.scores[move]
	plays := payoff.plays[move]
	heuristic := payoff.heuristic[move]

	exploitation := (1-alpha)*(float64(score)/float64(plays)) + alpha*heuristic
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/float64(plays))

	return exploitation + exploration
}

func calculateNodeHeuristic(node *Node, id fastGame.SnakeId) float64 {
	//closestFoodPath := FindNearestFood(node.Board, node.Ruleset, snake)
	health := float64(node.board.Healths[id])

	////foodScore := float64(1/len(closestFoodPath) + 1)
	////lengthScore := float64(len(snake.Body))

	////numOtherSnakes := 1
	////for _, s := range node.Board.Snakes {
	////if s.ID != snake.ID {
	////numOtherSnakes += 1
	////}
	////}
	////otherSnakeScore := float64(1 / numOtherSnakes)
	//a := 60.0
	//b := 8.0
	//foodDistance := float64(len(closestFoodPath))
	//foodScore := a * math.Atan(health-foodDistance/b)

	var otherSnakes []fastGame.SnakeId
	total := 0
	for id, health := range node.board.Healths {
		if id == fastGame.MeId && health > 0 {
			total += 100
		}
		if id != fastGame.MeId && health > 0 {
			otherSnakes = append(otherSnakes, id)
		}
	}
	snakeScore := 10 / (len(otherSnakes) + 1)
	healthScore := float64(health / 100)

	return float64(total) + float64(snakeScore) + healthScore
}
