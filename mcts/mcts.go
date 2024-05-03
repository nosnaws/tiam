package mcts

import (
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"

	bo "github.com/nosnaws/tiam/board"
)

type Node struct {
	board         *bo.FastBoard
	children      []*Node
	parent        *Node
	plays         int
	moveSet       map[bo.SnakeId]bo.SnakeMove
	possibleMoves map[bo.SnakeId][]bo.SnakeMove
	payoffs       map[bo.SnakeId]Payoff
	depth         int
}

type Payoff struct {
	plays     map[bo.Move]int
	scores    map[bo.Move]int
	heuristic map[bo.Move]float64
}

type SnakeScore struct {
	id        bo.SnakeId
	value     int
	move      bo.Move
	heuristic float64
}

func addAttributes(txn *newrelic.Transaction, root *Node, selected bo.SnakeMove, maxDepth int) {
	if txn != nil {
		txn.AddAttribute("totalPlays", root.plays)
		txn.AddAttribute("selectedMove", selected.Dir)
		txn.AddAttribute("maxDepth", maxDepth)
	}
}

func MCTS(board *bo.FastBoard, txn *newrelic.Transaction) bo.SnakeMove {
	segment := txn.StartSegment("MCTS")
	defer segment.End()
	fakeMoveSet := make(map[bo.SnakeId]bo.SnakeMove)
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
			node := selectNode(root)

			child := expandNode(node)

			score := simulateNode(child)
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

func selectNode(node *Node) *Node {
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node))
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

func simulateNode(node *Node) map[bo.SnakeId]SnakeScore {

	ns := node.board.Clone()

	//RandomRollout(&ns)
	StrategicRollout(&ns)

	//nodeHeuristic := calculateNodeHeuristic(node, bo.MeId)
	nodeHeuristic := float64(0)

	scores := make(map[bo.SnakeId]SnakeScore, len(ns.Heads))
	for id := range ns.Lengths {
		snakeHeuristic := nodeHeuristic
		if id != bo.MeId {
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

		scores[id] = score
	}

	return scores
}

func backpropagate(node *Node, scores map[bo.SnakeId]SnakeScore) {
	if node == nil {
		return
	}

	pastMovesWithScore := make(map[bo.SnakeId]SnakeScore)
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
	productOfMoves := bo.GetCartesianProductOfMoves(*node.board)

	var children []*Node
	for _, moveSet := range productOfMoves {
		cs := node.board.Clone()
		cs.AdvanceBoard(movesToMap(moveSet))

		moves := make(map[bo.SnakeId]bo.SnakeMove)
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

func movesToMap(moves []bo.SnakeMove) map[bo.SnakeId]bo.Move {
	m := make(map[bo.SnakeId]bo.Move, len(moves))
	for _, move := range moves {
		m[move.Id] = move.Dir
	}
	return m
}

func createNode(moveSet map[bo.SnakeId]bo.SnakeMove, board *bo.FastBoard) *Node {
	possibleMoves := make(map[bo.SnakeId][]bo.SnakeMove)
	payoffs := make(map[bo.SnakeId]Payoff)
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

func createPayoff(moves []bo.SnakeMove) Payoff {
	plays := make(map[bo.Move]int, len(moves))
	scores := make(map[bo.Move]int, len(moves))
	heuristic := make(map[bo.Move]float64, len(moves))

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

func selectFinalMove(node *Node) bo.SnakeMove {
	moves := node.possibleMoves[bo.MeId]
	sort.Slice(moves, func(a, b int) bool {
		return moveSecureness(node, bo.MeId, moves[a]) > moveSecureness(node, bo.MeId, moves[b])
	})

	return moves[0]
}

func moveSecureness(node *Node, player bo.SnakeId, move bo.SnakeMove) float64 {
	//score := float64(node.payoffs[player].scores[move.Dir])
	plays := float64(node.payoffs[player].plays[move.Dir])

	return plays
}

func bestUTC(node *Node) *Node {
	var moveSet []bo.SnakeMove
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

func isStateEqual(a []bo.SnakeMove, b map[bo.SnakeId]bo.SnakeMove) bool {
	equal := true
	for _, m := range a {
		if m.Dir != b[m.Id].Dir {
			equal = false
		}
	}

	return equal
}

func bestMoveUTC(node *Node, id bo.SnakeId) bo.SnakeMove {
	moves := node.possibleMoves[id]
	//sort.Slice(moves, func(a, b int) bool {
	//return calculateUCB(node, id, moves[a].Dir) > calculateUCB(node, id, moves[b].Dir)
	//})
	sort.Slice(moves, func(a, b int) bool {
		return calculateUCBTuned(node, id, moves[a].Dir) > calculateUCBTuned(node, id, moves[b].Dir)
	})
	//sort.Slice(moves, func(a, b int) bool {
	//return calculateUCBMinimal(node, id, moves[a].Dir) > calculateUCBMinimal(node, id, moves[b].Dir)
	//})

	return moves[0]
}

func calculateUCB(node *Node, id bo.SnakeId, move bo.Move) float64 {
	payoff := node.payoffs[id]
	explorationConstant := math.Sqrt(2)
	//alpha := float64(0.1)

	numParentSims := float64(node.plays)
	score := float64(payoff.scores[move])
	plays := float64(payoff.plays[move])
	//heuristic := payoff.heuristic[move]

	//exploitation := (1-alpha)*(score/plays) + alpha*heuristic
	exploitation := score / plays
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/plays)

	return exploitation + exploration
}

func calculateUCBMinimal(node *Node, id bo.SnakeId, move bo.Move) float64 {
	payoff := node.payoffs[id]

	score := float64(payoff.scores[move])
	plays := float64(payoff.plays[move])

	exploitation := score / plays
	exploration := 2 / plays

	return exploitation + exploration
}

func calculateUCBTuned(node *Node, id bo.SnakeId, move bo.Move) float64 {
	payoff := node.payoffs[id]
	//explorationConstant := math.Sqrt(2)
	//alpha := float64(0.1)

	numParentSims := float64(node.plays)
	score := float64(payoff.scores[move])
	plays := float64(payoff.plays[move])
	//heuristic := payoff.heuristic[move]
	variance := math.Pow(score, 2) / plays
	mean := nodeMean(node, id)
	v := variance - mean + math.Sqrt((2*math.Log(numParentSims))/plays)

	//exploitation := (1-alpha)*(score/plays) + alpha*heuristic
	exploitation := score / plays
	exploration := math.Sqrt((math.Log(numParentSims) / plays) * math.Min(1/4, v))

	//fmt.Println("ucb ", exploitation+exploration)
	return exploitation + exploration
}

func nodeMean(node *Node, id bo.SnakeId) float64 {
	payoff := node.payoffs[id]
	totalScore := 0
	totalSquared := 0.0
	for _, s := range payoff.scores {
		totalScore += s
		totalSquared += math.Pow(float64(s), 2)
	}
	mean := float64(totalScore / node.plays)
	return mean
}

func nodeVariance(node *Node, id bo.SnakeId) float64 {
	payoff := node.payoffs[id]
	totalScore := 0
	totalSquared := 0.0
	for _, s := range payoff.scores {
		totalScore += s
		totalSquared += math.Pow(float64(s), 2)
	}
	mean := float64(totalScore / node.plays)

	variance := (totalSquared / float64(node.plays)) - math.Pow(mean, 2)

	return variance
}

func nodeStdDev(node *Node, id bo.SnakeId) float64 {
	payoff := node.payoffs[id]
	totalScore := 0
	for _, s := range payoff.scores {
		totalScore += s
	}
	mean := float64(totalScore / node.plays)

	variance := 0.0
	for _, s := range payoff.scores {
		variance += math.Pow(float64(s)-mean, 2)
	}

	return math.Sqrt(variance / float64(node.plays))
}

func calculateNodeHeuristic(node *Node, id bo.SnakeId) float64 {
	//closestFoodPath := FindNearestFood(node.Board, node.Ruleset, snake)
	//health := float64(node.board.Healths[id])

	////foodScore := float64(1/len(closestFoodPath) + 1)
	////lengthScore := float64(len(snake.Body))

	////otherSnakeScore := float64(1 / numOtherSnakes)
	//a := 60.0
	//b := 8.0
	//foodDistance := float64(len(closestFoodPath))
	//foodScore := a * math.Atan(health-foodDistance/b)
	if !node.board.IsSnakeAlive(id) {
		return -1.0
	}

	//var otherSnakes []fastGame.SnakeId
	//otherSnakes := 0
	//for sId, health := range node.board.Healths {
	//if sId != id && health > 0 {
	//otherSnakes += 1
	//}
	//}
	//snakeScore := 1 / (otherSnakes + 1)
	//healthScore := 0.01 * float64(health/100)
	//lengthScore := 0.1 * float64(node.board.Lengths[id])
	vorRes := bo.Voronoi(node.board, id)
	voronoi := 0.01 * float64(vorRes.Score[id])
	total := voronoi

	return 1 / (1 + math.Pow(math.E, -total))
}
