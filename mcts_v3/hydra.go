package mctsv3

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/nosnaws/tiam/bitboard"
)

// setup workers to run the MCTS algorithm, return the results
// no memory is shared, move requests are sent to the workers via channels
//
// each request and pass along the context so that they can all be signalled to return their current result
// at the end, the main routine can aggregate the results to make a decision

// single request channel per worker
// one response channel

type requestChan chan workRequest
type responseChan chan rewards

// creates new workers each request, easier for now
func GetNextMove(board *bitboard.BitBoard, config MCTSConfig, me string) bitboard.Dir {
	workers := 8
	responses := make(responseChan, workers)
	duration, _ := time.ParseDuration("400ms")
	timeout, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		//b := board.Clone()
		go func(id int) {
			responses <- worker(timeout, id, board, config, me)
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(responses)

	totalRewards := make(rewards)

	for res := range responses {
		for resMove, resReward := range res {
			if _, ok := totalRewards[resMove]; !ok {
				totalRewards[resMove] = reward{}
			}

			r := totalRewards[resMove]
			r.plays += resReward.plays
			r.tac[0] += resReward.tac[0]
			r.tac[1] += resReward.tac[1]
			r.tac[2] += resReward.tac[2]
			totalRewards[resMove] = r
		}
	}

	fmt.Println("AGGREGATE", totalRewards)

	return bestMoveByTactic(totalRewards, 0)
}

func worker(ctx context.Context, id int, board *bitboard.BitBoard, config MCTSConfig, me string) rewards {
	root := CreateTree(board, me)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return MCTSWorker(
		ctx,
		rand,
		root,
		config,
	)
}

// controller
type Hydra struct {
	workers      []requestChan
	responseChan responseChan
}

type workRequest struct {
	duration time.Duration
	config   MCTSConfig
	board    *bitboard.BitBoard
	me       string
}

// worker
type HydraHead struct {
	requestChan  requestChan
	responseChan responseChan
}

//func CreateHydraHead(resChan responseChan) requestChan {
//requests := make(requestChan)
//head := HydraHead{
//responseChan: resChan,
//requestChan:  requests,
//}

//go head.spawnHead()

//return requests
//}

//func (hh *HydraHead) spawnHead() {
//for {
//select {
//case req := <-hh.requestChan:
//root := CreateTree(req.board, req.me)
//hh.responseChan <- MCTSWorker(
//root,
//req.duration,
//req.config.MinimumSims,
//req.config.SimLimit,
//req.config.TreeDepthLimit,
//)
//}

//}
//}
