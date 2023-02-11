package bitboard

import (
	"math/rand"

	"github.com/nosnaws/tiam/battlesnake"
)

func (bb *BitBoard) SpawnHazardsRoyale() {
	if bb.hazardSpawnTime < 1 {
		return
	}

	if bb.turn < bb.hazardSpawnTime {
		return
	}

	numShrinks := bb.turn / bb.hazardSpawnTime
	minX, maxX := 0, bb.width-1
	minY, maxY := 0, bb.height-1

	for i := 0; i < numShrinks; i++ {
		switch rand.Intn(4) {
		case 0:
			minX += 1
		case 1:
			maxX -= 1
		case 2:
			minY += 1
		case 3:
			maxY -= 1
		}
	}

	for x := 0; x < bb.width; x++ {
		for y := 0; y < bb.height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				// add hazard
				index := getIndex(battlesnake.Coord{X: x, Y: y}, bb.width)
				bb.hazards = bb.hazards.SetBit(index, 1)
			}
		}
	}
}
