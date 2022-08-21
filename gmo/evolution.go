package gmo

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"
)

// [0] = exploration constant
// [1] = alpha
// [2] = vorornoi weighting
// [3] = food/health weighting a
// [4] = food/health weighting b
// [5] = biggest snake reward
type genotype [6]float64

type canditate struct {
	geno      []float64
	name      string
	imageName string
}

type unnaturalSelection struct {
	currentPop           []canditate
	currentBest          canditate
	mutationConstant     float64
	mutationProbability  float64
	crossoverProbability float64
	roundLength          int
	populationSize       int
	numParents           int
	numGenerations       int
	scores               map[string]int
}

var randomSnake = canditate{
	geno:      []float64{0, 0, 0, 0, 0, 0},
	name:      "random_snake",
	imageName: "random",
}

var eaterSnake = canditate{
	geno:      []float64{0, 0, 0, 0, 0, 0},
	name:      "eater_snake",
	imageName: "eater",
}

type canditateWithScore struct {
	c canditate
	s int
}

func CreateEvolution(gens, popSize, numParents, roundLength int, mutConst, mutProb, crossProb float64) *unnaturalSelection {
	return &unnaturalSelection{
		populationSize:       popSize,
		numParents:           numParents,
		roundLength:          roundLength,
		mutationConstant:     mutConst,
		mutationProbability:  mutProb,
		crossoverProbability: crossProb,
		numGenerations:       gens,
	}
}

func (u *unnaturalSelection) Evolve(ctx context.Context) {
	// initialize the population
	initialPop := u.initializePopulation()
	u.currentPop = initialPop
	fmt.Println("initial pop ", initialPop)

	var overAllBest canditate

	for i := 0; i < u.numGenerations; i++ {
		fmt.Println("Running generation ", i)
		best, parents, scores := u.selection(ctx)

		writeCSV(i+1, u.currentPop, scores)
		if i == u.numGenerations-1 {
			overAllBest = best
			break
		}

		fmt.Println("Selected parent pool ", parents)
		fmt.Println("Adding best candidate to next gen ", best)
		nextgen := []canditate{best}

		for len(nextgen) < u.populationSize {
			pair := u.rouleteWheelSelection(parents, scores, 2)
			fmt.Println("Selected parents for breeding ", pair)

			newPair := u.crossover(pair[0], pair[1])

			fmt.Println("Adding child to next gen ", newPair[0])
			nextgen = append(nextgen, newPair[0])

			// avoid going over the pop size
			if len(nextgen) < u.populationSize {
				fmt.Println("Adding child to next gen ", newPair[1])
				nextgen = append(nextgen, newPair[1])
			}
		}

		mutatedPop := []canditate{}
		// potentially mutate
		for _, cand := range nextgen {
			mutatedPop = append(mutatedPop, u.mutate(cand))
		}

		u.currentPop = mutatedPop

	}

	fmt.Println("APEX", overAllBest)
}

func (u *unnaturalSelection) initializePopulation() []canditate {
	rand.Seed(time.Now().UnixNano())
	pop := []canditate{}

	for i := 0; i < u.populationSize; i++ {
		cand := canditate{
			name: genRandomName(),
			geno: []float64{
				rand.Float64() * float64(rand.Intn(5)),  // Exploration
				rand.Float64(),                          // alpha
				rand.Float64(),                          // voronoi
				rand.Float64(),                          // food a
				rand.Float64(),                          // food b
				rand.Float64() * float64(rand.Intn(10)), // big snake reward
			},
			imageName: "tiam",
		}
		pop = append(pop, cand)
	}

	return pop
}

func (u *unnaturalSelection) selection(ctx context.Context) (canditate, []canditate, map[string]int) {

	// Get the fitness of the population
	candidateSets := createGroups(u.currentPop)
	scores := make(map[string]int)
	for _, candidateSet := range candidateSets {
		fmt.Println("running group ", candidateSet)
		results := u.fitness(ctx, candidateSet)

		scores[candidateSet[0].name] = results[0]
		scores[candidateSet[1].name] = results[1]
	}

	// grab best
	best := u.currentPop[0]
	for _, c := range u.currentPop {
		if scores[c.name] > scores[best.name] {
			best = c
		}
	}
	// Select for reproduction

	selected := u.rouleteWheelSelection(u.currentPop, scores, u.numParents)

	return best, selected, scores
}

func (u *unnaturalSelection) rouleteWheelSelection(parents []canditate, scores map[string]int, n int) []canditate {
	rand.Seed(time.Now().UnixNano())

	totalScore := 0.0
	for _, p := range parents {
		totalScore += float64(scores[p.name])
	}

	sort.Slice(parents, func(i, j int) bool {
		return scores[parents[i].name] > scores[parents[j].name]
	})

	accScore := make(map[string]float64)
	acc := 0.0
	for _, cs := range parents {
		newScore := (float64(scores[cs.name]) / totalScore)
		accScore[cs.name] = newScore + acc
		acc += newScore
	}

	selected := []canditate{}
	for i := 0; i < n; i++ {
		randVal := rand.Float64()

		for _, cs := range parents {
			if accScore[cs.name] >= randVal {
				selected = append(selected, cs)
				break
			}
		}

	}

	return selected
}

func (u *unnaturalSelection) fitness(ctx context.Context, candidates []canditate) [2]int {
	can1 := candidates[0]
	can2 := candidates[1]
	deployment := createDeployment()

	// add candidates to deployment
	deployment.addContainer(can1.name, can1.imageName, genotypeToEnvVariables(can1.geno))
	deployment.addContainer(can2.name, can2.imageName, genotypeToEnvVariables(can2.geno))

	// add other snakes to deployment
	deployment.addContainer(randomSnake.name, randomSnake.imageName, []string{})
	deployment.addContainer(eaterSnake.name, eaterSnake.imageName, []string{})

	// deploy the snakes
	deployment.run(ctx)
	defer deployment.stopAndRemoveContainers(ctx)

	player1 := player{
		name: can1.name,
		url: fmt.Sprintf("http://localhost:%s",
			deployment.containers[can1.name].port),
	}
	player2 := player{
		name: can2.name,
		url: fmt.Sprintf("http://localhost:%s",
			deployment.containers[can2.name].port),
	}

	extra1 := player{
		name: randomSnake.name,
		url: fmt.Sprintf("http://localhost:%s",
			deployment.containers[randomSnake.name].port),
	}
	extra2 := player{
		name: eaterSnake.name,
		url: fmt.Sprintf("http://localhost:%s",
			deployment.containers[eaterSnake.name].port),
	}

	game := game{
		players: []player{player1, player2, extra1, extra2},
		mapName: "standard",
		height:  11,
		width:   11,
	}

	// run the games
	scores := [2]int{}
	for i := 0; i < u.roundLength; i++ {
		fmt.Println("Running game ", i)
		result := runGame(game, true)

		if result == player1.name {
			scores[0] += 1
		} else if result == player2.name {
			scores[1] += 1
		}
	}

	return scores
}

func (u *unnaturalSelection) crossover(canA, canB canditate) []canditate {
	rand.Seed(time.Now().UnixNano())
	// Return the original candidates if crossover is not performed
	if rand.Float64() > u.crossoverProbability {
		fmt.Println("No crossover, keeping genotype from parents")
		return []canditate{
			{
				name:      genRandomName(),
				imageName: canA.imageName,
				geno:      canA.geno,
			},
			{
				name:      genRandomName(),
				imageName: canB.imageName,
				geno:      canB.geno,
			},
		}
	}

	fmt.Println("Crossover ", canA.name, canB.name)
	crossOverPoint := rand.Intn(6)

	partA1, partA2 := canA.geno[:crossOverPoint], canA.geno[crossOverPoint:]
	partB1, partB2 := canB.geno[:crossOverPoint], canB.geno[crossOverPoint:]

	offspring1 := canditate{
		name:      genRandomName(),
		imageName: canA.imageName,
	}
	offspring1.geno = append(offspring1.geno, partA1...)
	offspring1.geno = append(offspring1.geno, partB2...)

	offspring2 := canditate{
		name:      genRandomName(),
		imageName: canA.imageName,
	}
	offspring2.geno = append(offspring2.geno, partB1...)
	offspring2.geno = append(offspring2.geno, partA2...)

	return []canditate{offspring1, offspring2}
}

func (u *unnaturalSelection) mutate(cand canditate) canditate {
	rand.Seed(time.Now().UnixNano())
	// Return the original candidate if mutation is not performed
	if rand.Float64() > u.mutationProbability {
		fmt.Println("No mutation, keeping genotype for ", cand)
		return cand
	}

	fmt.Println("Mutating ", cand.name)
	mutationPoint := rand.Intn(6)
	mutationSign := rand.Intn(2)

	if mutationSign == 0 {
		cand.geno[mutationPoint] = cand.geno[mutationPoint] + u.mutationConstant
	} else {
		cand.geno[mutationPoint] = cand.geno[mutationPoint] - u.mutationConstant
	}

	return cand
}

func createGroups(pop []canditate) [][]canditate {
	if len(pop)%2 != 0 {
		log.Fatal("Population size is not divisible by 2!")
	}

	cands := shuffle(pop)
	var pairs [][]canditate
	for i := 0; i < len(cands); i += 2 {
		pairs = append(pairs, []canditate{cands[i], cands[i+1]})
	}

	return pairs
}

func shuffle(cands []canditate) []canditate {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]canditate, len(cands))
	perm := r.Perm(len(cands))
	for i, randIndex := range perm {
		ret[i] = cands[randIndex]
	}
	return ret
}

func genotypeToEnvVariables(g []float64) []string {
	exploration := fmt.Sprintf("EXPLORATION_CONSTANT=%f", g[0])
	alpha := fmt.Sprintf("ALPHA_CONSTANT=%f", g[1])
	voronoi := fmt.Sprintf("VORONOI_CONSTANT=%f", g[2])
	foodA := fmt.Sprintf("FOOD_CONSTANT_A=%f", g[3])
	foodB := fmt.Sprintf("FOOD_CONSTANT_B=%f", g[3])
	biggest := fmt.Sprintf("BIG_SNAKE_CONSTANT=%f", g[3])

	return []string{exploration, alpha, voronoi, foodA, foodB, biggest}
}
