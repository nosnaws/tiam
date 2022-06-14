package main

import (
	"fmt"
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
}

type Payoff struct {
	Plays  map[string]int
	Scores map[string]int
}

type MoveScore struct {
	Direction string
	Plays     int
	Value     int
}

type SnakeScore struct {
	ID    string
	Value int
	Move  string
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

	duration, err := time.ParseDuration("200ms")
	if err != nil {
		panic("could not parse duration")
	}

	start := time.Now()

	for time.Since(start).Milliseconds() < duration.Milliseconds() {
		node := selectNode(root)
		child := expandNode(node)
		score := simulateNode(child)
		child.Plays += 1
		backpropagate(node, score)
	}

	fmt.Println("# ROOT #")
	fmt.Println("youID", youId)
	printNode(root)
	fmt.Println(root.PossibleMoves)
	fmt.Println("# Children #")
	for _, child := range root.Children {
		printNode(child)
	}

	bestMove := selectFinalMove(root)
	fmt.Println("# Selected #")
	fmt.Println(bestMove)
	addAttributes(txn, root, bestMove)
	return bestMove

	//fmt.Println("Could not find move, going left")
	//return rules.SnakeMove{ID: root.YouId, Move: "left"}
}

func printDepth(node *Node, acc int) {
	if node == nil {
		fmt.Println("Depth: ", acc)
		return
	}

	printDepth(node.Parent, acc+1)
}

func selectNode(node *Node) *Node {
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node))
}

func printNode(node *Node) {
	fmt.Println("Score", node.Payoffs)
	fmt.Println("Move", node.MoveSet)
}

func isGameOver(board *rules.BoardState, ruleset rules.Ruleset) bool {
	isGameOver, err := ruleset.IsGameOver(board)
	if err != nil {
		fmt.Println(board)
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

	var scores []SnakeScore
	for _, snake := range ns.Snakes {
		score := SnakeScore{ID: snake.ID, Value: 0, Move: node.MoveSet[snake.ID].Move}
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
		pastMovesWithScore = append(pastMovesWithScore, SnakeScore{ID: sc.ID, Value: sc.Value, Move: node.MoveSet[sc.ID].Move})
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
			fmt.Println(cs)
			fmt.Println(moveSet)
			panic("could not create next board state")
		}

		moves := make(map[string]rules.SnakeMove)
		for _, m := range moveSet {
			moves[m.ID] = m
		}

		childNode := createNode(node.YouId, moves, ns, node.Ruleset)
		childNode.Parent = node

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

	for _, m := range moves {
		plays[m.Move] = 0
		scores[m.Move] = 0
	}
	return Payoff{Plays: plays, Scores: scores}
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
		return calculateUTC(node, snakeId, moves[a].Move) > calculateUTC(node, snakeId, moves[b].Move)
	})

	return moves[0]
}

func calculateUTC(node *Node, snakeId string, move string) float64 {
	payoff := node.Payoffs[snakeId]
	explorationConstant := math.Sqrt(2)
	numParentSims := float64(node.Plays)
	score := payoff.Scores[move]
	plays := payoff.Plays[move]

	exploitation := float64(score) / float64(plays)
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/float64(plays))

	return exploitation + exploration
}
