package main

import (
	"fmt"

	"github.com/BattlesnakeOfficial/rules"
)

type Pair struct {
	snakeHead rules.Point
	point     rules.Point
}

func voronoiScore(node *Node) int {
	scores := voronoi(node.Board, node.Ruleset.Name() == "wrapped", 5)
	myScore := 0
	enemyTotal := 0
	for _, snake := range node.Board.Snakes {
		head := snake.Body[0]
		if snake.ID == node.YouId {
			myScore += scores[head]
		} else {
			enemyTotal += scores[head]
		}
	}
	return myScore
}

func voronoi(board *rules.BoardState, isWrapped bool, max int) map[rules.Point]int {
	depth := 0
	counts := make(map[rules.Point]int)
	visited := make(map[rules.Point]rules.Point)
	var q []Pair
	depthMark := Pair{snakeHead: rules.Point{X: -1, Y: -1}, point: rules.Point{X: -1, Y: -1}}
	mark := rules.Point{X: -1, Y: -1}

	for _, snake := range board.Snakes {
		head := snake.Body[0]
		p := Pair{snakeHead: head, point: head}
		q = append(q, p)
		visited[head] = head
		counts[head] = 0
	}
	q = append(q, depthMark)

	var p Pair
	for len(q) > 0 {
		fmt.Println(q)
		fmt.Println(counts)
		p, q = q[0], q[1:]
		fmt.Println("Exploring: ", p)

		if p == depthMark {
			fmt.Println("Found depthMark")
			depth += 1

			q = append(q, depthMark)

			if q[0] == depthMark || depth == max {
				fmt.Println("Done")
				break
			}
		} else {
			for _, neighbor := range GetEdges(p.point, board, isWrapped) {
				fmt.Println("Neighbor: ", neighbor)
				visitedFrom, ok := visited[neighbor]

				if ok {
					fmt.Println("Has been visited")
					if visitedFrom != mark && visitedFrom != p.snakeHead {
						fmt.Println("Has been visited by: ", visitedFrom)
						counts[visitedFrom] -= 1
						visited[neighbor] = mark
					}
				} else {
					fmt.Println("Has not been visited, adding to queue")
					counts[p.snakeHead] += 1
					visited[neighbor] = p.snakeHead
					q = append(q, Pair{snakeHead: p.snakeHead, point: neighbor})
				}
			}
		}
	}

	return counts
}

func BFSSpace(board *rules.BoardState, snake rules.Snake, isWrapped bool, max int) int {
	depth := 0
	count := 0
	visited := make(map[rules.Point]bool)
	var q []rules.Point
	depthMark := rules.Point{X: -1, Y: -1}

	head := snake.Body[0]
	q = append(q, head)
	visited[head] = true

	q = append(q, depthMark)

	var p rules.Point
	for len(q) > 0 {
		fmt.Println(q)
		p, q = q[0], q[1:]
		fmt.Println("Exploring: ", p)

		if p == depthMark {
			fmt.Println("Found depthMark")
			depth += 1

			q = append(q, depthMark)

			if q[0] == depthMark || depth == max {
				fmt.Println("Done")
				break
			}
		} else {
			for _, neighbor := range GetEdges(p, board, isWrapped) {
				fmt.Println("Neighbor: ", neighbor)

				if visited[neighbor] == false {
					visited[neighbor] = true
					count += 1
					q = append(q, neighbor)
				}

			}
		}
	}

	return count
}
