package main

import (
	"context"

	"github.com/nosnaws/tiam/gmo"
)

func main() {
	ctx := context.Background()
	generations := 10
	numParents := 4
	populationSize := 10
	roundSize := 1
	mutationCon := 0.2
	mutationProb := 0.1
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
