package main

import (
	"fmt"

	"github.com/BattlesnakeOfficial/rules"
)

func BuildRuleset(game GameState) rules.Ruleset {
	name := game.Game.Ruleset.Name
	settings := game.Game.Ruleset.Settings

	return rules.NewRulesetBuilder().WithParams(map[string]string{

		"name":              name,
		"foodSpawnChance":   "0",
		"minimumFood":       "0",
		"shrinkEveryNTurns": "0",
		"damagePerTurn":     fmt.Sprint(settings.HazardDamagePerTurn),
	}).Ruleset()
}

func BuildBoard(game GameState) rules.BoardState {
	snakes := game.Board.Snakes

	return rules.BoardState{
		Turn:    int32(game.Turn),
		Height:  int32(game.Board.Height),
		Width:   int32(game.Board.Width),
		Food:    coordsToPoints(game.Board.Food),
		Hazards: coordsToPoints(game.Board.Hazards),
		Snakes:  battleSnakesToSnakes(snakes),
	}
}

func battleSnakesToSnakes(battleSnakes []Battlesnake) []rules.Snake {
	var snakes []rules.Snake
	for _, s := range battleSnakes {
		snakes = append(snakes, battleSnakeToSnake(s))
	}

	return snakes
}

func FindNearestFood(board *rules.BoardState, ruleset rules.Ruleset, snake rules.Snake) []rules.Point {
	isWrapped := ruleset.Name() == "wrapped"
	return breadthFirstSearchFood(board, snake, isWrapped)
}

type BFSEdge struct {
	Point  rules.Point
	Parent *BFSEdge
}

func breadthFirstSearchFood(board *rules.BoardState, snake rules.Snake, isWrapped bool) []rules.Point {
	current := snake.Body[0]
	var q []rules.Point
	explored := make(map[rules.Point]*BFSEdge)
	explored[current] = &BFSEdge{Point: current}
	q = append(q, current)

	var v rules.Point
	for len(q) > 0 {
		v, q = q[0], q[1:]

		if isPointInSlice(v, board.Food) {
			return constructBFSPath(explored[v])
		}

		for _, edge := range GetEdges(v, board, isWrapped) {
			if explored[edge] == nil && !isPointInSlice(edge, board.Hazards) {
				explored[edge] = &BFSEdge{Point: edge, Parent: explored[v]}
				q = append(q, edge)
			}
		}
	}

	return []rules.Point{}
}

func constructBFSPath(edge *BFSEdge) []rules.Point {
	if edge == nil {
		return []rules.Point{}
	}

	return append(constructBFSPath(edge.Parent), edge.Point)
}

func isPointInSlice(p rules.Point, s []rules.Point) bool {
	for _, current := range s {
		if current == p {
			return true
		}
	}
	return false
}

// http://prtamil.github.io/posts/cartesian-product-go/
func GetCartesianProductOfMoves(board *rules.BoardState, ruleset rules.Ruleset) [][]rules.SnakeMove {
	var allMoves [][]rules.SnakeMove
	for _, snake := range board.Snakes {
		moves := GetSnakeMoves(snake, ruleset, *board)
		allMoves = append(allMoves, moves)
	}

	var temp [][]rules.SnakeMove
	for _, a := range allMoves[0] {
		temp = append(temp, []rules.SnakeMove{a})
	}

	for i := 1; i < len(allMoves); i++ {
		temp = cartesianProduct(temp, allMoves[i])
	}

	return temp
}

func cartesianProduct(movesA [][]rules.SnakeMove, movesB []rules.SnakeMove) [][]rules.SnakeMove {
	var result [][]rules.SnakeMove
	for _, a := range movesA {
		for _, b := range movesB {
			var temp []rules.SnakeMove
			for _, m := range a {
				temp = append(temp, m)
			}

			temp = append(temp, b)
			result = append(result, temp)
		}
	}

	return result
}

func getSnake(id string, board *rules.BoardState) *rules.Snake {
	for _, snake := range board.Snakes {
		if snake.ID == id {
			return &snake
		}
	}
	return nil
}

func battleSnakeToSnake(bs Battlesnake) rules.Snake {
	return rules.Snake{
		ID:     bs.ID,
		Health: bs.Health,
		Body:   coordsToPoints(bs.Body),
	}
}

func coordsToPoints(coords []Coord) []rules.Point {
	var points []rules.Point
	for _, c := range coords {
		points = append(points, coordToPoint(c))
	}

	return points
}

func coordToPoint(c Coord) rules.Point {
	return rules.Point{X: int32(c.X), Y: int32(c.Y)}
}

var possibleMoves = []string{rules.MoveUp, rules.MoveDown, rules.MoveLeft, rules.MoveRight}

func GetEdges(head rules.Point, board *rules.BoardState, isWrapped bool) []rules.Point {
	edges := []rules.Point{{X: head.X + 1, Y: head.Y}, {X: head.X - 1, Y: head.Y}, {X: head.X, Y: head.Y + 1}, {X: head.X, Y: head.Y - 1}}
	edges = updateForWrapped(edges, board.Height, board.Width, isWrapped)
	edges = filterInBounds(edges, board.Height, board.Width)
	return edges
}

func updateForWrapped(p []rules.Point, h int32, w int32, isWrapped bool) []rules.Point {
	if !isWrapped {
		return p
	}

	var wrappedPoints []rules.Point
	for _, point := range p {
		wrappedPoints = append(wrappedPoints, rules.Point{X: point.X % w, Y: point.Y % h})
	}

	return wrappedPoints
}

func filterInBounds(p []rules.Point, h int32, w int32) []rules.Point {
	var inBounds []rules.Point
	for _, point := range p {
		if point.X >= w || point.Y >= h || point.X < 0 || point.Y < 0 {
			continue
		}

		inBounds = append(inBounds, point)
	}

	return inBounds
}

func GetSnakeMoves(snake rules.Snake, ruleset rules.Ruleset, board rules.BoardState) []rules.SnakeMove {
	head := snake.Body[0]
	neck := snake.Body[1]

	safeMoves := []string{}

	// avoid walls
	if ruleset.Name() != "wrapped" {
		for _, move := range possibleMoves {
			isValid := false
			if move == rules.MoveUp && head.Y+1 < board.Height {
				isValid = true
			}
			if move == rules.MoveDown && head.Y-1 >= 0 {
				isValid = true
			}
			if move == rules.MoveRight && head.X+1 < board.Width {
				isValid = true
			}
			if move == rules.MoveLeft && head.X-1 >= 0 {
				isValid = true
			}

			if isValid {
				safeMoves = append(safeMoves, move)
			}

		}
	} else {
		safeMoves = possibleMoves
	}

	nonNeckMoves := []string{}
	// avoid the neck
	for _, move := range safeMoves {
		isValid := false
		if move == rules.MoveUp && head.Y+1 != neck.Y {
			isValid = true
		}
		if move == rules.MoveDown && head.Y-1 != neck.Y {
			isValid = true
		}
		if move == rules.MoveRight && head.X+1 != neck.X {
			isValid = true
		}
		if move == rules.MoveLeft && head.X-1 != neck.X {
			isValid = true
		}

		if isValid {
			nonNeckMoves = append(nonNeckMoves, move)
		}

	}

	snakeMoves := []rules.SnakeMove{}
	for _, m := range nonNeckMoves {
		// skip hazards
		if isPointInSlice(dirToPoint(head, m), board.Hazards) {
			continue
		}
		snakeMoves = append(snakeMoves, createSnakeMove(snake.ID, m))
	}

	return snakeMoves
}

func dirToPoint(head rules.Point, dir string) rules.Point {
	if dir == rules.MoveLeft {
		return rules.Point{X: head.X - 1, Y: head.Y}
	}
	if dir == rules.MoveRight {
		return rules.Point{X: head.X + 1, Y: head.Y}
	}
	if dir == rules.MoveUp {
		return rules.Point{X: head.X, Y: head.Y + 1}
	}
	return rules.Point{X: head.X, Y: head.Y - 1}
}

func createSnakeMove(id string, move string) rules.SnakeMove {
	return rules.SnakeMove{ID: id, Move: move}
}
