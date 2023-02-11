package bitboard

import (
	"math/rand"

	"github.com/nosnaws/tiam/battlesnake"
	"github.com/shabbyrobe/go-num"
)

func (bb *BitBoard) SpawnFood() {
	// need to find the total number of food
	if bb.totalFood < bb.minimumFood {
		bb.placeFoodRandomly(bb.minimumFood - bb.totalFood)
		return
	}

	if bb.foodSpawnChance > 0 && int(rand.Intn(100)) < bb.foodSpawnChance {
		bb.placeFoodRandomly(1)
	}
}

func (bb *BitBoard) placeFoodRandomly(num int) {
	for i := 0; i < num; i++ {
		unoccupiedIndices := bb.getUnoccupiedIndices()

		if len(unoccupiedIndices) > 0 {
			newFood := unoccupiedIndices[rand.Intn(len(unoccupiedIndices))]
			bb.food = bb.food.SetBit(newFood, 1)
			bb.totalFood += 1
		}
	}
}

func (bb *BitBoard) getUnoccupiedIndices() []int {
	indices := []int{}

	for i := 0; i < bb.width*bb.height; i++ {
		if bb.empty.Bit(i) > 0 && bb.food.Bit(i) == 0 {
			indices = append(indices, i)
		}
	}
	return indices
}

func (bb *BitBoard) placeFoodFixed() {
	centerCoord := battlesnake.Coord{X: (bb.width - 1) / 2, Y: (bb.height - 1) / 2}

	for _, snake := range bb.Snakes {
		headCoord := indexToPoint(snake.GetHeadIndex(), bb.width)
		possibleFoodLocations := []battlesnake.Coord{
			{X: headCoord.X - 1, Y: headCoord.Y - 1},
			{X: headCoord.X - 1, Y: headCoord.Y + 1},
			{X: headCoord.X + 1, Y: headCoord.Y - 1},
			{X: headCoord.X + 1, Y: headCoord.Y + 1},
		}

		availableFoodLocations := []battlesnake.Coord{}
		for _, p := range possibleFoodLocations {
			if p.X < 0 || p.X > bb.width-1 || p.Y < 0 || p.Y > bb.height-1 {
				continue
			}

			coordIndex := getIndex(p, bb.width)

			if centerCoord == p {
				continue
			}

			coordBoard := num.U128From16(0)
			coordBoard = coordBoard.SetBit(coordIndex, 1)

			// skip if there is already food
			if coordBoard.And(bb.food).BitLen() > 0 {
				continue
			}

			isAwayFromCenter := false
			if p.X < headCoord.X && headCoord.X < centerCoord.X {
				isAwayFromCenter = true
			} else if centerCoord.X < headCoord.X && headCoord.X < p.X {
				isAwayFromCenter = true
			} else if p.Y < headCoord.Y && headCoord.Y < centerCoord.Y {
				isAwayFromCenter = true
			} else if centerCoord.Y < headCoord.Y && headCoord.Y < p.Y {
				isAwayFromCenter = true
			}
			if !isAwayFromCenter {
				continue
			}

			// no food in corners
			if (p.X == 0 || p.X == (bb.width-1)) && (p.Y == 0 || p.Y == (bb.height-1)) {
				continue
			}

			availableFoodLocations = append(availableFoodLocations, p)
		}

		if len(availableFoodLocations) <= 0 {
			return
		}

		placedFood := availableFoodLocations[rand.Intn(len(availableFoodLocations))]
		bb.food = bb.food.SetBit(getIndex(placedFood, bb.width), 1)
	}

	isCenterOccupied := true
	centerBoard := num.U128From16(0)
	centerBoard = centerBoard.SetBit(getIndex(centerCoord, bb.width), 1)

	if centerBoard.And(bb.empty).BitLen() == 0 {
		isCenterOccupied = false
	}

	if isCenterOccupied {
		return
	}

	bb.food = bb.food.SetBit(getIndex(centerCoord, bb.width), 1)
}
