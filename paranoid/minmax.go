package paranoid

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/moveset"
)

type MinMaxScore struct {
	move  moveset.MoveSet
	score float64
}

// original scores
//total += foodScore * 2
//total += killScore
//total += areaScore * 0.8
//total += length * 0.5
//total += health * 0.1
func GetMove(board *bitboard.BitBoard, maxId string, maxDepth int) bitboard.Dir {
	weights := bitboard.BasicStateWeights{
		Food:   2,
		Aggr:   1,
		Area:   0.8,
		Length: 0.5,
		Health: 0.1,
	}
	result, _ := negamax(context.TODO(), weights, board, maxId, 0, 0, maxDepth, -math.MaxFloat64, math.MaxFloat64, true)

	fmt.Println("MOVE", moveset.Dir(result.move), result.score)

	return bitboard.MoveSetToDir(result.move)
}

func GetMoveArena(board *bitboard.BitBoard, maxId string, maxDepth int, weights bitboard.BasicStateWeights) bitboard.SnakeMoveSet {
	ctx := context.TODO()
	result, _ := negamax(ctx, weights, board, maxId, 0, 0, maxDepth, -math.MaxFloat64, math.MaxFloat64, true)

	return bitboard.SnakeMoveSet{
		Id:  maxId,
		Set: result.move,
	}
}

func GetMoveID(board *bitboard.BitBoard, maxId string) bitboard.Dir {
	weights := bitboard.BasicStateWeights{
		Food:   4.17007,
		Aggr:   5.97133,
		Area:   6.268947,
		Length: -0.352937,
		Health: 1.473651,
		Tail:   2.456908,
	}

	d, _ := time.ParseDuration("450ms")
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	result := IDNegamax(ctx, weights, board, maxId)

	log.Println("MOVE", moveset.ToDirs(result.move), result.score)

	return bitboard.MoveSetToDir(result.move)
}

func IDNegamax(ctx context.Context, weights bitboard.BasicStateWeights, board *bitboard.BitBoard, maxId string) MinMaxScore {
	currentDepth := 0
	currentMove := MinMaxScore{}

negamaxLoop:
	for {
		select {
		case <-ctx.Done():
			break negamaxLoop
		default:
		}

		currentDepth += 2

		if currentDepth > 100 {
			log.Println("Over 100 cutting off")
			break negamaxLoop
		}

		log.Println("Running depth", currentDepth)
		result, ignore := negamax(ctx, weights, board, maxId, 0, 0, currentDepth, -math.MaxFloat64, math.MaxFloat64, true)
		log.Println("Result", result, ignore)

		if ignore {
			continue
		}

		currentMove = result
	}

	return currentMove
}

// TODO: move ordering
func negamax(ctx context.Context, weights bitboard.BasicStateWeights, board *bitboard.BitBoard, maxId string, maxMove moveset.MoveSet, depth, maxDepth int, alpha, beta float64, isMax bool) (MinMaxScore, bool) {

	select {
	case <-ctx.Done():
		return MinMaxScore{
			move:  maxMove,
			score: 0,
		}, true
	default:
	}

	if depth >= maxDepth || isTerminalNode(board, maxId) {
		return MinMaxScore{
			move:  maxMove,
			score: board.BasicStateScore(maxId, depth, weights),
		}, false
	}

	if isMax {

		moves := moveset.Split(board.GetMoves(maxId).Set)
		ignoreRes := false

		value := -math.MaxFloat64
		move := maxMove

		for _, posMove := range moves {
			negResult, ignore := negamax(ctx, weights, board, maxId, posMove, depth+1, maxDepth, -beta, -alpha, false)
			ignoreRes = ignore
			negResult.score = -negResult.score

			if negResult.score > value {
				value = negResult.score
				move = negResult.move
			}

			if value > alpha {
				alpha = value
			}

			if alpha >= beta {
				break
			}
		}

		return MinMaxScore{
			move:  move,
			score: value,
		}, ignoreRes
	} else {
		// opponents

		moves := withOppMoves(board, maxId, maxMove)
		ignoreRes := false

		value := -math.MaxFloat64

		for _, moveSet := range moves {
			ns := board.Clone()
			ns.AdvanceTurn(moveSet)

			negResult, ignore := negamax(ctx, weights, ns, maxId, maxMove, depth+1, maxDepth, -beta, -alpha, true)
			ignoreRes = ignore
			negResult.score = -negResult.score

			if negResult.score > value {
				value = negResult.score
			}

			if value > alpha {
				alpha = value
			}

			if alpha >= beta {
				break
			}
		}

		return MinMaxScore{
			move:  maxMove,
			score: value,
		}, ignoreRes
	}
}

func withOppMoves(board *bitboard.BitBoard, maxId string, maxMove moveset.MoveSet) [][]bitboard.SnakeMoveSet {

	allMoves := [][]bitboard.SnakeMoveSet{}
	allMoves = append(allMoves, []bitboard.SnakeMoveSet{{
		Id:  maxId,
		Set: maxMove,
	}})

	for id := range board.Snakes {
		if id == maxId {
			continue
		}

		moves := moveset.Split(board.GetMoves(id).Set)
		moveSets := []bitboard.SnakeMoveSet{}
		for _, m := range moves {
			moveSets = append(moveSets, bitboard.SnakeMoveSet{
				Id:  id,
				Set: m,
			})
		}

		allMoves = bitboard.CartesianProduct(allMoves, moveSets)
	}

	return allMoves
}

func isTerminalNode(board *bitboard.BitBoard, maxId string) bool {
	if !board.IsSnakeAlive(maxId) || moveset.IsEmpty(board.GetMovesNoDefault(maxId).Set) {
		return true
	}

	if board.IsGameOver() {
		return true
	}

	return false
}
