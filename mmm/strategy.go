package mmm

import (
	"math"

	b "github.com/nosnaws/tiam/board"
)

type StrategyFn func(*b.FastBoard, b.SnakeId, b.SnakeId, int) float64

func multiStrat(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	//evoStrat := GMOConfig{
	//FoodWeightA:    4.166351,
	//FoodWeightB:    3.785205,
	//VoronoiWeightA: 1.900364,
	//VoronoiWeightB: 1.150082,
	//LengthWeight:   2.334582,
	//}
	//return StrategyGMO(board, maxId, minId, depth, evoStrat)
	return StrategyHungry(board, maxId, minId, depth)
}

func StrategyV1(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	//health := float64(board.Healths[id])
	log := getLogger()

	if isDeadOrOut(board, maxId) {
		log.debugAlways(depth, "We lose!")
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		log.debugAlways(depth, "We win!")
		return float64(10000 * (depth + 1))
	}

	minMoves := board.GetMovesForSnakeNoDefault(minId)
	minMIndexs := []uint16{}
	for _, mm := range minMoves {
		minMIndexs = append(minMIndexs, b.IndexInDirection(mm.Dir, board.Heads[minId], board.Width, board.Height, board.IsWrapped))
	}
	//voronoi := b.Voronoi(board, id)
	maxTail := board.GetSnakeTail(maxId)
	ff, foodDepth, tailDepth := floodfill(board, int(board.Heads[maxId]), int(board.Lengths[maxId]*2), []uint16{maxTail})
	minFF, _, _ := floodfill(board, int(board.Heads[minId]), int(board.Lengths[minId]*2), []uint16{})

	total := 0.0
	//total += float64(board.Lengths[maxId]) * 0.02
	//total += float64(voronoi.Score[id]) * 0.01
	//total += float64(voronoi.FoodDepth[id]) * 0.01
	total += float64(ff) / 2
	total += -(float64(minFF) / 2)

	if ff < int(board.Lengths[maxId]) {
		// less space available than our length is bad
		// but not quite as bad as losing the game
		total += float64(-50 * (depth + 1))
	}

	isLargerSnake := board.Lengths[maxId] > board.Lengths[minId]
	if isLargerSnake {
		enemyDepth := findBF(board, int(board.Heads[maxId]), func(i int) bool {
			for _, mm := range minMoves {
				mIndex := board.MoveToIndex(mm)
				if int(mIndex) == i {
					return true
				}
			}
			return false
		})
		total += float64(-enemyDepth)
	}

	if foodDepth != -1 {
		total += float64(-foodDepth) * 2
	} else {
		food := findBF(board, int(board.Heads[maxId]), func(i int) bool {
			return board.IsTileFood(uint16(i))
		})

		if food != -1 {
			total += float64(-food) * 2
		}
	}
	// if we are near a hazard, that's not the best
	if hNeighbors := board.GetHazardNeighbors(board.Heads[maxId]); len(hNeighbors) > 0 {
		total += float64(-20 * len(hNeighbors))
	}

	// moving toward our tail is usually good
	if tailDepth != -1 {
		total += float64(-tailDepth * 2)
	} else {
		total += float64(-50 * (depth + 1))
	}

	return total
}

func StrategyV2(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	//health := float64(board.Healths[id])
	log := getLogger()

	if isDeadOrOut(board, maxId) {
		log.debugAlways(depth, "We lose!")
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		log.debugAlways(depth, "We win!")
		return float64(10000 * (depth + 1))
	}

	//minMoves := board.GetMovesForSnakeNoDefault(minId)
	//minMIndexs := []uint16{}
	//for _, mm := range minMoves {
	//minMIndexs = append(minMIndexs, b.IndexInDirection(mm.Dir, board.Heads[minId], board.Width, board.Height, board.IsWrapped))
	//}
	voronoi := b.Voronoi(board, maxId)

	total := 0.0
	//total += float64(board.Lengths[maxId]) * 0.02
	otherSnakes := getOtherSnakeIds(board, maxId)

	if len(otherSnakes) == 1 {
		isLargerSnake := board.Lengths[maxId] > board.Lengths[minId]
		v := float64(voronoi.Score[maxId])

		fd := float64(voronoi.FoodDepth[maxId])
		if fd != -1 {
			if !isLargerSnake {
				fd = fd * 2
			}
			total += -fd
		}

		total += v
	} else {
		ff, foodDepth, _ := floodfill(board, int(board.Heads[maxId]), int(board.Lengths[maxId]*2), []uint16{})

		if foodDepth == -1 {
			foodDepth = findBF(board, int(board.Heads[maxId]), func(i int) bool {
				return board.IsTileFood(uint16(i))
			})
		}

		if foodDepth != -1 {
			total += float64(-foodDepth) * 2
		}

		total += float64(ff)
	}

	maxTail := board.GetSnakeTail(maxId)
	tailDepth := findBF(board, int(board.Heads[maxId]), func(i int) bool {
		if int(maxTail) == i {
			return true
		}
		return false
	})

	// moving toward our tail is usually good
	if tailDepth != -1 {
		total += float64(-tailDepth * 2)
	} else {
		total += float64(-50 * (depth + 1))
	}

	return total
}

func StrategyV3(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	if isDeadOrOut(board, maxId) {
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		return float64(10000 * (depth + 1))
	}

	health := float64(board.Healths[maxId])
	//length := float64(board.Lengths[maxId])
	// TODO: modify voronoi to only run for 2 snakes
	voronoi := b.Voronoi(board, maxId)

	total := 0.0

	maxV := float64(voronoi.Score[maxId])
	minV := float64(voronoi.Score[minId])

	fd := float64(voronoi.FoodDepth[maxId])
	if fd != -1 {
		total += 12 * (health - fd) / 20
	}

	total += (maxV - minV)

	maxTail := board.GetSnakeTail(maxId)
	tailDepth := findBF(board, int(board.Heads[maxId]), func(i int) bool {
		if int(maxTail) == i {
			return true
		}
		return false
	})

	// not being able to loop back is bad
	if tailDepth == -1 {
		total += float64(-500 * (depth + 1))
	}

	return total
}

func StrategyV4(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	if isDeadOrOut(board, maxId) {
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		return float64(10000 * (depth + 1))
	}

	//health := float64(board.Healths[maxId])
	//length := float64(board.Lengths[maxId])
	// TODO: modify voronoi to only run for 2 snakes
	voronoi := b.Voronoi(board, maxId)

	total := 0.0

	//total += length * 2

	maxV := float64(voronoi.Score[maxId])
	//minV := float64(voronoi.Score[minId])

	fd := float64(voronoi.FoodDepth[maxId])
	if fd != -1 {
		fd = -fd * 4
		//fmt.Println("ADDING FOOD DEPTH", fd)
		total += fd
	}

	//fmt.Println("ADDING VORONOI", (maxV))
	total += 0.5 * (maxV)

	//articulations := getArticulationPoints(board, board.Heads[maxId])
	//for _, a := range articulations {
	//if voronoi.Territory[a] == maxId {
	//total += 20
	//}
	//}

	//fmt.Println("SCORE", total)
	return total
}

func StrategyV5(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	if isDeadOrOut(board, maxId) {
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		return float64(10000 * (depth + 1))
	}

	//health := float64(board.Healths[maxId])
	//length := float64(board.Lengths[maxId])
	// TODO: modify voronoi to only run for 2 snakes
	//voronoi := b.Voronoi(board, maxId)
	ff, foodDepth, _ := floodfill(board, int(board.Heads[maxId]), 32, []uint16{})
	//minFF, _, _ := floodfill(board, int(board.Heads[minId]), 32, []uint16{})

	total := 0.0

	//total += length * 2

	//maxV := float64(voronoi.Score[maxId])
	//minV := float64(voronoi.Score[minId])

	fd := float64(foodDepth)
	if fd != -1 {
		fd = -fd * 4
		//fmt.Println("ADDING FOOD DEPTH", fd)
		total += fd
	}

	total += float64(ff)

	//fmt.Println("ADDING VORONOI", (maxV))

	//articulations := getArticulationPoints(board, board.Heads[maxId])
	//for _, a := range articulations {
	//if voronoi.Territory[a] == maxId {
	//total += 20
	//}
	//}

	//fmt.Println("SCORE", total)
	return total
}

type GMOConfig struct {
	FoodWeightA    float64
	FoodWeightB    float64
	VoronoiWeightA float64
	VoronoiWeightB float64
	LengthWeight   float64
}

func CreateGMOStrategy(gmo GMOConfig) StrategyFn {
	return func(board *b.FastBoard, maxId b.SnakeId, minId b.SnakeId, depth int) float64 {
		return StrategyGMO(board, maxId, minId, depth, gmo)
	}
}

func StrategyGMO(board *b.FastBoard, maxId, minId b.SnakeId, depth int, gmo GMOConfig) float64 {
	if isDeadOrOut(board, maxId) {
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		return float64(10000 * (depth + 1))
	}

	health := float64(board.Healths[maxId])
	maxLength := float64(board.Lengths[maxId])
	minLength := float64(board.Lengths[minId])
	voronoi := b.Voronoi(board, maxId)

	total := 0.0

	maxV := float64(voronoi.Score[maxId])
	minV := float64(voronoi.Score[minId])

	fd := float64(voronoi.FoodDepth[maxId])
	if fd == -1 {
		fd = 0
	}

	total += gmo.LengthWeight * (maxLength - minLength)
	total += gmo.FoodWeightA * math.Atan(health-fd) / gmo.FoodWeightB
	total += gmo.VoronoiWeightA * (maxV - minV) / gmo.VoronoiWeightB

	return total
}

func StrategyGMO2(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	if isDeadOrOut(board, maxId) {
		return float64(-10000 * (depth + 1))
	}

	if isDeadOrOut(board, minId) || board.IsGameOver() {
		return float64(10000 * (depth + 1))
	}

	voronoi := b.Voronoi(board, maxId)
	space, _, minDis := floodfill(board, int(board.Heads[maxId]), 20, []uint16{board.Heads[minId]})
	arts := getArticulationPoints(board, board.Heads[maxId])

	total := 0.0

	maxV := float64(voronoi.Score[maxId])

	fd := float64(voronoi.FoodDepth[maxId])
	if fd == -1 {
		fd = 0
	}

	for _, a := range arts {
		if voronoi.Territory[a] == maxId {
			total += 10
		}
	}
	total += float64(minDis)
	total += 1 * float64(space)
	total += 2 * -fd
	total += 1 * maxV

	return total
}

func StrategyHungry(board *b.FastBoard, maxId, minId b.SnakeId, depth int) float64 {
	maxOut := isDeadOrOut(board, maxId)
	minOut := isDeadOrOut(board, minId)
	if maxOut && minOut {
		return float64(-5000 * (depth + 1))
	}

	if maxOut {
		return float64(-10000 * (depth + 1))
	}

	if minOut || board.IsGameOver() {
		return float64(10000 * (depth + 1))
	}

	minMoves := board.GetMovesForSnake(minId)
	minMoveI := []uint16{}
	for _, mm := range minMoves {
		minMoveI = append(minMoveI, board.MoveToIndex(mm))
	}
	length := int(board.Lengths[maxId])
	//minLength := int(board.Lengths[minId])
	head := board.Heads[maxId]
	//minHead := board.Heads[minId]
	voronoi := b.Voronoi(board, maxId)
	maxV := float64(voronoi.Score[maxId])
	//minV := float64(voronoi.Score[minId])
	//if length < 4 {
	//maxV /= 10
	//}
	space, _, _ := floodfill(board, int(board.Heads[maxId]), 100, []uint16{board.Heads[minId]})
	//space := floodfillSafe(board, int(head), length*2, minMoveI)
	//space := dfSpaceFill(board, int(head), length*2)
	//arts := getArticulationPoints(board, board.Heads[maxId])
	health := float64(board.Healths[maxId])
	allFood := board.GetAllFood()
	totalFood := float64(len(allFood))
	tail := board.GetSnakeTail(maxId)
	tailDistance := findBFUnsafe(board, int(head), func(i int) bool {
		return uint16(i) == tail
	})
	isLargestSnake := true
	for id, l := range board.Lengths {
		if id != maxId && length-int(l) < 1 {
			isLargestSnake = false
		}
	}

	//lengthDiff := length - minLength
	shouldEat := !isLargestSnake || health < 40
	//isFirstEat := board.Lengths[maxId] < 4
	isTrap := space < length

	total := 0.0

	if isTrap {
		total += -float64(500 * (depth + 1))
	}

	if shouldEat {
		foodDistances := findAllBF(board, int(head), func(i int) bool {
			return board.IsTileFood(uint16(i))
		})
		//isFoodInCell := len(foodDistances) > 0
		spaceScore := 0.0 * float64(space)
		vorScore := 0.8 * maxV

		//fmt.Println("voronoi", vorScore)
		//fmt.Println("food", foodDistances)
		//fmt.Println("space", spaceScore)
		foodScore := 0.0
		for _, d := range foodDistances {
			if voronoi.Territory[d.Index] == maxId {
				foodScore += float64(d.Depth)
			}
		}
		//if length-minLength < 2 {
		//foodScore /= 4
		//}

		totalFoodScore := 0.4 * totalFood
		foodTotal := math.Atan(4/(foodScore+totalFoodScore)) * 6.3
		total += foodTotal * (vorScore + spaceScore)
	} else {
		spaceScore := 0.0 * float64(space)
		vorScore := 0.8 * maxV

		//fmt.Println("voronoi", maxV)
		//fmt.Println("voronoiMin", minV)
		//fmt.Println("space", spaceScore)
		//if length-minLength < 2 {
		//foodScore /= 4
		//}

		//minDistTotal := math.Atan(4/minDistScore) * 6.3
		// mutliplying by 10 so that the score is always larger
		// than a shouldEat score
		total += 10 * (vorScore + spaceScore)
	}
	if tailDistance != -1 {
		tailScore := math.Atan(4/(float64(tailDistance+1))) * 8
		//fmt.Println("tail", tailScore)
		total += tailScore
	}

	//moves := board.GetMovesForSnakeNoDefault(maxId)
	//total += float64(len(moves)) * 2
	//healthWeight := health / 2
	//total += foodTotal / (healthWeight)

	//fd := float64(voronoi.FoodDepth[maxId])
	//if fd == -1 {
	//fd = 0
	//}
	//total += -fd

	//artsScore := 0.0
	//for _, a := range arts {
	//if voronoi.Territory[a] == maxId && length > 4 {
	//artsScore += 1
	//}
	//}
	//total += float64(minDis)
	//total += 1 * float64(space)
	//total += 2 * -fd
	//total = (total + 1) * (maxV )
	//if tailDistance == -1 {
	//fmt.Println("COULD NOT FIND TAIL")
	//total += float64(-500 * (depth + 1))
	//}

	//fmt.Println("total", total)
	return total
}
