package brain

import (
	g "github.com/nosnaws/tiam/game"
)

type mmNode struct {
	board g.Board
	move  g.SnakeMove
}

type moveMatrix map[g.Move][][]g.SnakeMove
type abMatrix map[g.Move][][]int32

func minmax(gamestate g.FastBoard, move g.SnakeMove, depth int, alpha int32, beta int32, maxPlayer g.SnakeId) {
	ns := gamestate.Clone()
	if ns.IsGameOver() {
		// return move and score
	} else {

		// move matrix
		// M = max player, N = min players
		// M(i), [n1, n2, n3], [n4, n5, n6] ...

		//moveMatrix := createMoveMatrix(ns, maxPlayer)
		//alphaMatrix := createABMatrix(moveMatrix, alpha)
		//betaMatrix := createABMatrix(moveMatrix, beta)

		//for maxM, minMs := range moveMatrix {
		//for i, minM := range minMs {

		//}
		//}

	}

}

func createMoveMatrix(gs g.FastBoard, id g.SnakeId) moveMatrix {
	moveMatrix := make(moveMatrix, 4)
	movesCombos := g.GetCartesianProductOfMoves(gs)
	for _, moves := range movesCombos {
		var otherMoves []g.SnakeMove
		var maxMove g.SnakeMove
		for _, move := range moves {
			if move.Id == id {
				maxMove = move
			} else {
				otherMoves = append(otherMoves, move)
			}
		}
		moveMatrix[maxMove.Dir] = append(moveMatrix[maxMove.Dir], otherMoves)
	}
	return moveMatrix
}

func createABMatrix(mm moveMatrix, val int32) abMatrix {
	abMatrix := make(abMatrix, 4)
	for maxM, minMs := range mm {
		abMatrix[maxM] = make([][]int32, len(minMs))
		for i, minM := range minMs {
			abMatrix[maxM][i] = make([]int32, len(minM))
			for j := range minM {
				abMatrix[maxM][i][j] = val
			}
		}
	}
	return abMatrix
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
