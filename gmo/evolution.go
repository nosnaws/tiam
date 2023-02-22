package gmo

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nosnaws/tiam/arena"
	"github.com/nosnaws/tiam/arena/agents"
	bitboard "github.com/nosnaws/tiam/bitboard2"
	"github.com/nosnaws/tiam/paranoid"
)

// [0] = food/health weighting a
// [1] = food/health weighting b
// [2] = vorornoi weighting a
// [3] = vorornoi weighting b
// [4] = length weight
type genotype [5]float64

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
	geno:      []float64{0, 0, 0, 0, 0},
	name:      "random_snake",
	imageName: "random",
}

var eaterSnake = canditate{
	geno:      []float64{0, 0, 0, 0, 0},
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
	//initialPop := u.initializePopulation()
	initialPop := u.initializePopFromCSV("/Users/aswanson/dev/innovation/tiam/gen_60_feb21.csv")
	u.currentPop = initialPop
	fmt.Println("initial pop ", initialPop)

	var overAllBest canditate

	for i := 0; i < u.numGenerations; i++ {
		fmt.Println("Running generation ", i)
		best, scores := u.selection(ctx)

		writeCSV(i, u.currentPop, scores)
		if i == u.numGenerations-1 {
			overAllBest = best
			break
		}

		fmt.Println("Adding best candidate to next gen ", best)
		nextgen := []canditate{best}

		for len(nextgen) < u.populationSize {
			pair := u.rouleteWheelSelection(u.currentPop, scores, 2)
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

func randFloat(min, max int) float64 {
	return float64(min) + rand.Float64()*float64(max-min)
}

func (u *unnaturalSelection) initializePopulation() []canditate {
	pop := []canditate{}

	for i := 0; i < u.populationSize; i++ {
		cand := canditate{
			name: genRandomName(),
			geno: []float64{
				randFloat(0, 10), // food
				randFloat(0, 10), // aggresive
				randFloat(0, 10), // area
				randFloat(0, 10), // length
				randFloat(0, 10), // health
			},
			imageName: "huey",
		}
		pop = append(pop, cand)
	}

	return pop
}

func (u *unnaturalSelection) initializePopFromCSV(path string) []canditate {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Unable to read initialization file", path)
	}

	contents := string(file)
	rows := strings.Split(contents, "\n")
	rows = rows[1:]           // remove the header row
	rows = rows[:len(rows)-1] // remove the empty last line

	fmt.Println(len(rows))
	pop := []canditate{}
	for _, row := range rows {
		columns := strings.Split(row, ",")
		fmt.Println(columns, len(columns))
		// columns[1] is the fitness score

		g1, _ := strconv.ParseFloat(columns[2], 64)
		g2, _ := strconv.ParseFloat(columns[3], 64)
		g3, _ := strconv.ParseFloat(columns[4], 64)
		g4, _ := strconv.ParseFloat(columns[5], 64)
		g5, _ := strconv.ParseFloat(columns[6], 65)
		cand := canditate{
			name: columns[0],
			geno: []float64{
				g1,
				g2,
				g3,
				g4,
				g5,
			},
			imageName: "huey",
		}

		pop = append(pop, cand)
	}

	return pop
}

func (u *unnaturalSelection) selection(ctx context.Context) (canditate, map[string]int) {

	// Get the fitness of the population
	candidateSets := createGroups(u.currentPop)
	scores := make(map[string]int)

	for i := 0; i < len(candidateSets); i += 2 {
		ch1 := make(chan [2]int)
		ch2 := make(chan [2]int)

		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			log.Println("running group ", i, candidateSets[i])
			ch1 <- u.fitnessArena(ctx, candidateSets[i])
			wg.Done()
		}()
		go func() {
			log.Println("running group ", i+1, candidateSets[i+1])
			ch2 <- u.fitnessArena(ctx, candidateSets[i+1])
			wg.Done()
		}()

		results1 := <-ch1
		results2 := <-ch2
		wg.Wait()

		scores[candidateSets[i][0].name] = results1[0]
		scores[candidateSets[i][1].name] = results1[1]

		scores[candidateSets[i+1][0].name] = results2[0]
		scores[candidateSets[i+1][1].name] = results2[1]
	}

	// grab best
	best := u.currentPop[0]
	for _, c := range u.currentPop {
		if scores[c.name] > scores[best.name] {
			best = c
		}
	}

	return best, scores
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

func (u *unnaturalSelection) fitnessArena(ctx context.Context, candidates []canditate) [2]int {
	can1 := candidates[0]
	can2 := candidates[1]

	round := arena.Arena{
		Agents: []arena.Agent{
			{
				Id: can1.name,
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return paranoid.GetMoveArena(bb, s, 2, bitboard.BasicStateWeights{
						Food:   can1.geno[0],
						Aggr:   can1.geno[1],
						Area:   can1.geno[2],
						Length: can1.geno[3],
						Health: can1.geno[4],
					},
					)
				},
			},
			{
				Id: can2.name,
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return paranoid.GetMoveArena(bb, s, 2, bitboard.BasicStateWeights{
						Food:   can2.geno[0],
						Aggr:   can2.geno[1],
						Area:   can2.geno[2],
						Length: can2.geno[3],
						Health: can2.geno[4],
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
				Id: "random",
				GetMove: func(bb *bitboard.BitBoard, s string) bitboard.SnakeMoveSet {
					return agents.GetRandomMove(bb, s)
				},
			},
		},
		Rounds: u.roundLength,
	}

	round.Initialize()

	round.Run()

	scores := [2]int{}

	for _, result := range round.Results {
		if result.Winner == can1.name {
			scores[0] += 1
		} else if result.Winner == can2.name {
			scores[1] += 1
		}
	}

	log.Println("Finished group")

	return scores
}

func (u *unnaturalSelection) fitness(ctx context.Context, candidates []canditate) [2]int {
	can1 := candidates[0]
	can2 := candidates[1]
	deployment := createDeployment(8080)

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
		players:  []player{player1, player2, extra1, extra2},
		mapName:  "hz_islands_bridges",
		gameType: "wrapped",
		hzDamage: 100,
		height:   11,
		width:    11,
	}

	var wg sync.WaitGroup

	// run the games
	scores := [2]int{}
	var mutex = &sync.RWMutex{}
	for i := 0; i < u.roundLength; i++ {
		log.Println("Starting game ", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := runGame(game, false)

			mutex.Lock()
			if result == player1.name {
				scores[0] += 1
			} else if result == player2.name {
				scores[1] += 1
			}
			mutex.Unlock()
		}()
	}

	wg.Wait()
	log.Println("Finished group")

	return scores
}

func (u *unnaturalSelection) crossover(canA, canB canditate) []canditate {
	// Return the original candidates if crossover is not performed
	if rand.Float64() > u.crossoverProbability {
		log.Println("No crossover, keeping genotype from parents")
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

	log.Println("Crossover ", canA.name, canB.name)
	crossOverPoint := rand.Intn(5)

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
		log.Println("No mutation, keeping genotype for ", cand)
		return cand
	}

	log.Println("Mutating ", cand.name)
	mutationPoint := rand.Intn(5)
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
	foodA := fmt.Sprintf("FOOD_A=%f", g[0])
	foodB := fmt.Sprintf("FOOD_B=%f", g[1])
	vorA := fmt.Sprintf("VORONOI_A=%f", g[2])
	vorB := fmt.Sprintf("VORONOI_B=%f", g[3])
	length := fmt.Sprintf("LENGTH=%f", g[4])

	return []string{foodA, foodB, vorA, vorB, length}
}
