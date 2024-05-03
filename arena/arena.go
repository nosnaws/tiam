package arena

import (
	"fmt"
	"log"
	"math/rand"

	api "github.com/nosnaws/tiam/battlesnake"
	bitboard "github.com/nosnaws/tiam/bitboard2"
)

// Package for running lots of games and reporting results
// runs 11x11 boards
// supports standard and royale maps
// no support for wrapped

type AgentMoveFn func(*bitboard.BitBoard, string) bitboard.SnakeMoveSet

type Agent struct {
	Id      string
	GetMove AgentMoveFn
}

type RoundResult struct {
	Id         int
	Winner     string
	TotalTurns int
}

type Arena struct {
	Agents       []Agent
	Rounds       int
	Results      []RoundResult
	initialBoard *bitboard.BitBoard
}

var STARTING_TOP_LEFT = []api.Coord{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}
var STARTING_TOP_RIGHT = []api.Coord{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}
var STARTING_BOTTOM_LEFT = []api.Coord{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}
var STARTING_BOTTOM_RIGHT = []api.Coord{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}

func (a *Arena) Initialize() {
	a.initialBoard = initializeBoard(a.Agents)
}

func (a *Arena) Run() {
	log.Printf("Running rounds...")
	for i := 0; i < a.Rounds; i++ {
		result := a.runGame(i)

		a.Results = append(a.Results, result)
	}
	log.Printf("Done!")
}

func (a *Arena) runGame(id int) RoundResult {
	game := a.initialBoard.Clone()
	//d, _ := time.ParseDuration("500ms")

	for !game.IsGameOver() {
		moves := []bitboard.SnakeMoveSet{}

		for _, agent := range a.Agents {
			if game.IsSnakeAlive(agent.Id) {
				ns := game.Clone()
				moves = append(moves, agent.GetMove(ns, agent.Id))
			}
		}

		game.AdvanceTurn(moves)
		//game.Print()
		//time.Sleep(d)
	}

	return RoundResult{
		Winner:     determineWinner(game),
		Id:         id,
		TotalTurns: game.Turn,
	}
}

func determineWinner(board *bitboard.BitBoard) string {
	winner := "draw"

	for id, snake := range board.Snakes {
		if snake.IsAlive() {
			winner = id
		}
	}

	return winner
}

func initializeBoard(agents []Agent) *bitboard.BitBoard {
	if len(agents) > 4 {
		panic("Too many agents")
	}

	snakes := []api.Battlesnake{}

	boardPositions := [][]api.Coord{
		STARTING_BOTTOM_LEFT,
		STARTING_BOTTOM_RIGHT,
		STARTING_TOP_LEFT,
		STARTING_TOP_RIGHT,
	}

	rand.Shuffle(len(agents), func(i, j int) {
		agents[i], agents[j] = agents[j], agents[i]
	})

	for i, agent := range agents {
		snakes = append(snakes, api.Battlesnake{
			ID:     agent.Id,
			Health: 100,
			Body:   boardPositions[i],
		})

	}

	state := api.GameState{
		Turn: 0,
		Board: api.Board{
			Snakes: snakes,
			Height: 11,
			Width:  11,
			Food: []api.Coord{
				{X: 5, Y: 5},
			},
		},
		You: snakes[0],
	}
	board := bitboard.CreateBitBoard(state)
	board.ShouldSpawnFood = true

	return board
}

func printResults(results []RoundResult) {
	winners := make(map[string]int)
	total := len(results)

	for _, result := range results {
		winners[result.Winner] += 1
	}

	fmt.Println("RESULTS")

	for name, wins := range winners {
		fmt.Printf("%s: %d - %f\n", name, wins, float64(wins)/float64(total))
	}
}
