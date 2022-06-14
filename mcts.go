package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/BattlesnakeOfficial/rules"
)

type Node struct {
	youId    string
	score    int32
	plays    int32
	ruleset  rules.Ruleset
	board    *rules.BoardState
	children []*Node
	parent   *Node
	move     rules.SnakeMove
}

func MCTS(youId string, board *rules.BoardState, ruleset rules.Ruleset) rules.SnakeMove {
	fakeMove := rules.SnakeMove{ID: youId, Move: ""}
	root := createNode(youId, fakeMove, board, ruleset)
	root.children = createChildren(root)

	duration, err := time.ParseDuration("200ms")
	if err != nil {
		panic("could not parse duration")
	}

	start := time.Now()

	for time.Since(start).Milliseconds() < duration.Milliseconds() {
		node := selectNode(root)
		child := expandNode(node)
		score := simulateNode(child)
		backpropagate(child, score)
	}

	fmt.Println("# ROOT #")
	printNode(root)
	fmt.Println("# Children #")
	for _, child := range root.children {
		printNode(child)
	}

	bestChild := selectFinalMove(root)
	fmt.Println("# Selected #")
	printNode(bestChild)

	return bestChild.move
}

func printDepth(node *Node, acc int) {
	if node == nil {
		fmt.Println("Depth: ", acc)
		return
	}

	printDepth(node.parent, acc+1)
}

func selectNode(node *Node) *Node {
	if isLeafNode(node) {
		return node
	}

	return selectNode(bestUTC(node))
}

func printNode(node *Node) {
	fmt.Println("Plays", node.plays)
	fmt.Println("Score", node.score)
	fmt.Println("Move", node.move.Move)
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
	if isGameOver(node.board, node.ruleset) {
		return node
	}

	if len(node.children) == 0 {
		children := createChildren(node)

		node.children = children
	}

	return getRandomUnexploredChild(node)
}

func simulateNode(node *Node) int32 {
	ns := node.board.Clone()

	for isGameOver(ns, node.ruleset) == false {
		var allMoves []rules.SnakeMove
		for _, snake := range ns.Snakes {
			if snake.EliminatedCause != rules.NotEliminated {
				continue
			}

			moves := GetSnakeMoves(snake, node.ruleset, *ns)
			randomMove := moves[rand.Intn(len(moves))]
			allMoves = append(allMoves, randomMove)
		}

		ns, _ = node.ruleset.CreateNextBoardState(ns, allMoves)
	}

	you := getSnake(node.youId, ns)
	if you == nil {
		panic("cannot find self")
	}

	if you.EliminatedCause == rules.NotEliminated {
		return 1
	}

	return 0
}

func backpropagate(node *Node, score int32) {
	if node == nil {
		return
	}

	node.plays += 1
	node.score += score
	backpropagate(node.parent, score)
}

func isLeafNode(node *Node) bool {
	if len(node.children) == 0 {
		return true
	}

	for _, child := range node.children {
		if child.plays == 0 {
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

func getMove(id string, moves []rules.SnakeMove) *rules.SnakeMove {
	for _, move := range moves {
		if move.ID == id {
			return &move
		}
	}
	return nil
}

func createChildren(node *Node) []*Node {
	productOfMoves := GetCartesianProductOfMoves(node.board, node.ruleset)

	var children []*Node
	for _, moveSet := range productOfMoves {
		cs := node.board.Clone()
		ns, err := node.ruleset.CreateNextBoardState(cs, moveSet)
		if err != nil {
			fmt.Println(cs)
			fmt.Println(moveSet)
			panic("could not create next board state")
		}

		youMove := getMove(node.youId, moveSet)
		if youMove == nil {
			panic("could not find move")
		}

		childNode := createNode(node.youId, *youMove, ns, node.ruleset)
		childNode.parent = node

		children = append(children, childNode)
	}

	return children
}

func createNode(youId string, move rules.SnakeMove, board *rules.BoardState, rules rules.Ruleset) *Node {
	return &Node{youId: youId, move: move, board: board, ruleset: rules}
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

func selectFinalMove(node *Node) *Node {
	children := node.children
	sort.Slice(children, func(a, b int) bool {
		return children[a].plays > children[b].plays
	})

	return children[0]
}

func bestUTC(node *Node) *Node {
	children := node.children
	sort.Slice(children, func(a, b int) bool {
		return calculateUTC(children[a]) > calculateUTC(children[b])
	})

	return children[0]
}

func calculateUTC(node *Node) float64 {
	explorationConstant := math.Sqrt(2)
	numParentSims := float64(node.parent.plays)

	exploitation := float64(node.score) / float64(node.plays)
	exploration := explorationConstant * math.Sqrt(math.Log(numParentSims)/float64(node.plays))

	return exploitation + exploration
}
