package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Node struct {
	YouId         string
	Ruleset       rules.Ruleset
	Board         *rules.BoardState
	Children      []*Node
	Parent        *Node
	Plays         int
	MoveSet       map[string]rules.SnakeMove
	PossibleMoves map[string][]rules.SnakeMove
	Payoffs       map[string]Payoff
	Depth         int
}

type Payoff struct {
	Plays     map[string]int
	Scores    map[string]int
	Heuristic map[string]float64
}

type SnakeScore struct {
	ID        string
	Value     int
	Move      string
	Heuristic float64
}

func addAttributes(txn *newrelic.Transaction, root *Node, selected rules.SnakeMove) {
	txn.AddAttribute("totalPlays", root.Plays)
	txn.AddAttribute("selectedMove", selected.Move)
	for _, m := range root.PossibleMoves[root.YouId] {
		scoreKey := fmt.Sprintf("moves.scores.%s", m.Move)
		playsKey := fmt.Sprintf("moves.plays.%s", m.Move)
		txn.AddAttribute(scoreKey, root.Payoffs[root.YouId].Scores[m.Move])
		txn.AddAttribute(playsKey, root.Payoffs[root.YouId].Plays[m.Move])
	}
}

func MCTS(youId string, board *rules.BoardState, ruleset rules.Ruleset, txn *newrelic.Transaction) rules.SnakeMove {
	fakeMoveSet := make(map[string]rules.SnakeMove)
	root := createNode(youId, fakeMoveSet, board, ruleset)
	root.Children = createChildren(root)

	duration, err := time.ParseDuration("250ms")
	if err != nil {
		panic("could not parse duration")
	}

	//start := time.Now()

	//for time.Since(start).Milliseconds() < duration.Milliseconds() {
	//node := selectNode(root)
	//child := expandNode(node)
	//score := simulateNode(child)
	//child.Plays += 1
	//backpropagate(node, score)
	//}

	//defer txn.StartSegment("simulateNode").End()
loop:
	for timeout := time.After(duration); ; {
		select {
		case <-timeout:
			break loop
		default:
			t := txn.StartSegment("selectNode")
			node := selectNode(root)
			t.End()

			t = txn.StartSegment("expandNode")
			child := expandNode(node)
			t.End()

			t = txn.StartSegment("simulateNode")
			score := simulateNode(child)
			child.Plays += 1
			t.End()

			t = txn.StartSegment("backpropagate")
			backpropagate(node, score)
			t.End()
		}
	}

	log.Println("# ROOT #")
	log.Println("youID", youId)
	printNode(root)
	log.Println(root.PossibleMoves)
	log.Println("# Children #")
	for _, child := range root.Children {
		printNode(child)
	}

	bestMove := selectFinalMove(root)
	log.Println("# Selected #")
	log.Println(bestMove)
	addAttributes(txn, root, bestMove)
	return bestMove

	//log.Println("Could not find move, going left")
	//return rules.SnakeMove{ID: root.YouId, Move: "left"}
}

func selectNode(node *Node) *Node {
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node))
}

func printNode(node *Node) {
	log.Println("#############")
	log.Println("Depth", node.Depth)
	log.Println("Plays", node.Payoffs[node.YouId].Plays)
	log.Println("Scores", node.Payoffs[node.YouId].Scores)
	log.Println("Heuristics", node.Payoffs[node.YouId].Heuristic)
	log.Println("#############")
}

func isGameOver(board *rules.BoardState, ruleset rules.Ruleset) bool {
	isGameOver, err := ruleset.IsGameOver(board)
	if err != nil {
		log.Println(board)
		panic("tried to check if game was over")
	}

	return isGameOver
}

func expandNode(node *Node) *Node {
	if isGameOver(node.Board, node.Ruleset) {
		return node
	}

	if len(node.Children) == 0 {
		children := createChildren(node)

		node.Children = children
	}

	return getRandomUnexploredChild(node)
}

func simulateNode(node *Node) []SnakeScore {

	ns := node.Board.Clone()

	for isGameOver(ns, node.Ruleset) == false {
		var allMoves []rules.SnakeMove
		for _, snake := range ns.Snakes {
			if snake.EliminatedCause != rules.NotEliminated {
				continue
			}

			moves := GetSnakeMoves(snake, node.Ruleset, *ns)
			randomMove := moves[rand.Intn(len(moves))]
			allMoves = append(allMoves, randomMove)
		}

		ns, _ = node.Ruleset.CreateNextBoardState(ns, allMoves)
	}

	var nodeHeuristic float64
	for _, snake := range node.Board.Snakes {
		if snake.ID == node.YouId {
			nodeHeuristic = calculateNodeHeuristic(node, snake)
		}
	}

	var scores []SnakeScore
	for _, snake := range ns.Snakes {
		snakeHeuristic := nodeHeuristic
		if snake.ID != node.YouId {
			snakeHeuristic = -snakeHeuristic
		}

		score := SnakeScore{ID: snake.ID, Value: 0, Move: node.MoveSet[snake.ID].Move, Heuristic: snakeHeuristic}
		if snake.EliminatedCause == rules.NotEliminated {
			score.Value = 1
		}
		scores = append(scores, score)
	}

	return scores
}

func backpropagate(node *Node, scores []SnakeScore) {
	if node == nil {
		return
	}

	var pastMovesWithScore []SnakeScore
	node.Plays += 1
	for _, sc := range scores {
		node.Payoffs[sc.ID].Plays[sc.Move] += 1
		node.Payoffs[sc.ID].Scores[sc.Move] += sc.Value

		h := node.Payoffs[sc.ID].Heuristic[sc.Move]
		node.Payoffs[sc.ID].Heuristic[sc.Move] = math.Max(h, sc.Heuristic)

		pastMovesWithScore = append(pastMovesWithScore, SnakeScore{ID: sc.ID, Value: sc.Value, Move: node.MoveSet[sc.ID].Move, Heuristic: sc.Heuristic})
	}

	backpropagate(node.Parent, pastMovesWithScore)
}

func isLeafNode(node *Node) bool {
	if len(node.Children) == 0 {
		return true
	}

	for _, n := range node.Children {
		if n.Plays < 1 {
			return true
		}
	}

	return false
}

func getRandomUnexploredChild(node *Node) *Node {
	var unexplored []*Node
	for _, child := range node.Children {
		if child.Plays == 0 {
			unexplored = append(unexplored, child)
		}
	}

	return Shuffle(unexplored)[0]
}

func getMove(id string, moves []rules.SnakeMove) *rules.SnakeMove {
	for _, move := range moves {
		if move.ID == id {
			return &move
		}
	}
	return nil
}

func createChildren(node *Node) []*Node {
	productOfMoves := GetCartesianProductOfMoves(node.Board, node.Ruleset)

	var children []*Node
	for _, moveSet := range productOfMoves {
		cs := node.Board.Clone()
		ns, err := node.Ruleset.CreateNextBoardState(cs, moveSet)
		if err != nil {
			log.Println(cs)
			log.Println(moveSet)
			panic("could not create next board state")
		}

		moves := make(map[string]rules.SnakeMove)
		for _, m := range moveSet {
			moves[m.ID] = m
		}

		childNode := createNode(node.YouId, moves, ns, node.Ruleset)
		childNode.Parent = node
		childNode.Depth = childNode.Parent.Depth + 1

		children = append(children, childNode)
	}

	return children
}

func createNode(youId string, moveSet map[string]rules.SnakeMove, board *rules.BoardState, ruleset rules.Ruleset) *Node {
	possibleMoves := make(map[string][]rules.SnakeMove)
	payoffs := make(map[string]Payoff)
	for _, snake := range board.Snakes {
		moves := GetSnakeMoves(snake, ruleset, *board)
		possibleMoves[snake.ID] = moves
		payoffs[snake.ID] = createPayoff(moves)
	}

	return &Node{YouId: youId, PossibleMoves: possibleMoves, Board: board, Ruleset: ruleset, Payoffs: payoffs, MoveSet: moveSet}
}

func createPayoff(moves []rules.SnakeMove) Payoff {
	plays := make(map[string]int)
	scores := make(map[string]int)
	heuristics := make(map[string]float64)

	for _, m := range moves {
		plays[m.Move] = 0
		scores[m.Move] = 0
		heuristics[m.Move] = 0
	}
	return Payoff{Plays: plays, Scores: scores, Heuristic: heuristics}
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

func selectFinalMove(node *Node) rules.SnakeMove {
	moves := node.PossibleMoves[node.YouId]
	payoff := node.Payoffs[node.YouId]
	sort.Slice(moves, func(a, b int) bool {
		return payoff.Scores[moves[a].Move] > payoff.Scores[moves[b].Move]
	})

	return moves[0]
}

func bestUTC(node *Node) *Node {
	var moveSet []rules.SnakeMove
	for _, snake := range node.Board.Snakes {
		bestMove := bestMoveUTC(node, snake.ID)
		moveSet = append(moveSet, bestMove)
	}

	for _, child := range node.Children {
		if isStateEqual(moveSet, child.MoveSet) {
			return child
		}
	}

	return nil
}

func isStateEqual(a []rules.SnakeMove, b map[string]rules.SnakeMove) bool {
	equal := true
	for _, m := range a {
		if m.Move != b[m.ID].Move {
			equal = false
		}
	}

	return equal
}

func bestMoveUTC(node *Node, snakeId string) rules.SnakeMove {
	moves := node.PossibleMoves[snakeId]
	sort.Slice(moves, func(a, b int) bool {
		return calculateUCB(node, snakeId, moves[a].Move) > calculateUCB(node, snakeId, moves[b].Move)
	})

	return moves[0]
}

func calculateUCB(node *Node, snakeId string, move string) float64 {
	payoff := node.Payoffs[snakeId]
	explorationConstant := math.Sqrt(2)
	alpha := float64(0.1)

	numParentSims := float64(node.Plays)
	score := payoff.Scores[move]
	plays := payoff.Plays[move]
	heuristic := payoff.Heuristic[move]

	exploitation := (1-alpha)*(float64(score)/float64(plays)) + alpha*heuristic
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/float64(plays))

	return exploitation + exploration
}

func calculateNodeHeuristic(node *Node, snake rules.Snake) float64 {
	//closestFoodPath := FindNearestFood(node.Board, node.Ruleset, snake)
	health := float64(snake.Health)

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

	var otherSnakes []rules.Snake
	total := 0
	for _, snake := range node.Board.Snakes {
		if snake.ID == node.YouId && snake.EliminatedCause == rules.NotEliminated {
			total += 100
		}
		if snake.ID != node.YouId && snake.EliminatedCause == rules.NotEliminated {
			otherSnakes = append(otherSnakes, snake)
		}
	}
	snakeScore := 10 / (len(otherSnakes) + 1)
	healthScore := float64(health / 100)

	return float64(total) + float64(snakeScore) + healthScore
}
