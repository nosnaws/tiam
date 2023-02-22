package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/nosnaws/tiam/gmo"
)

func main() {
	rand.Seed(time.Now().UnixMicro())

	ctx := context.Background()
	generations := 60
	numParents := 20
	populationSize := 200
	roundSize := 20
	mutationCon := 0.8
	mutationProb := 0.2
	crossProb := 0.8

	evo := gmo.CreateEvolution(
		generations,
		populationSize,
		numParents,
		roundSize,
		mutationCon,
		mutationProb,
		crossProb,
	)
	evo.Evolve(ctx)
}
