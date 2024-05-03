package arena

import (
	"math/rand"
	"testing"
	"time"

	"github.com/nosnaws/tiam/arena/agents"
	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/paranoid"
)

func TestArenaRound(t *testing.T) {
	rand.Seed(time.Now().UnixMicro())

	arena := Arena{
		Agents: []Agent{
			{
				Id: "heuy",
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return paranoid.GetMoveArena(bb, s, 2, bitboard.BasicStateWeights{
						Food:   4.165581,
						Aggr:   12.623248,
						Area:   7.454674,
						Length: 8.788976,
						Health: 3.698813,
					},
					)
				},
			},
			{
				Id: "eater",
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return agents.GetEaterMove(bb, s)
				},
			},
			{
				Id: "3",
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return agents.GetRandomMove(bb, s)
				},
			},
			{
				Id: "4",
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return agents.GetRandomMove(bb, s)
				},
			},
		},
		Rounds: 1000,
	}

	arena.initialBoard = initializeBoard(arena.Agents)

	arena.Run()

	printResults(arena.Results)
	if len(arena.Results) != arena.Rounds {
		panic("did not run 2 rounds")
	}
}
