package gmo

import (
	"context"
	"fmt"
	"testing"
)

func TestFitness(t *testing.T) {
	t.Skip()
	can1 := canditate{
		geno:      []float64{0, 0, 0, 0},
		name:      "test_tiam_1",
		imageName: "tiam",
	}
	can2 := canditate{
		geno:      []float64{0, 0, 0, 0},
		name:      "test_tiam_2",
		imageName: "tiam",
	}

	evo := unnaturalSelection{
		roundLength: 5,
	}

	res := evo.fitness(context.TODO(), []canditate{can1, can2})

	fmt.Println(res)
}

func TestEvolution(t *testing.T) {
	t.Skip()

	evo := CreateEvolution(20, 10, 5, 2, 0.2, 0.1, 0.8)
	evo.Evolve(context.TODO())

}
