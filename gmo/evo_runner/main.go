package main

import (
	"context"

	"github.com/nosnaws/tiam/gmo"
)

func main() {
	ctx := context.Background()
	generations := 60
	numParents := 10
	populationSize := 100
	roundSize := 10
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
