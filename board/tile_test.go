package board

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tile := CreateEmptyTile()

	if tile.IsEmpty() != true {
		panic("Not empty!")
	}

	tile = CreateHeadTile(0, 0)
	if tile.IsEmpty() != false {
		panic("Should not be empty!")
	}
}

func TestIsFood(t *testing.T) {
	tile := CreateEmptyTile()

	if tile.IsFood() != false {
		panic("Should not be food!")
	}

	tile.SetFood()
	if tile.IsFood() != true {
		panic("Should be food!")
	}
}

func TestHazardInteractions(t *testing.T) {
	tile := CreateEmptyTile()

	if tile.IsHazard() != false {
		panic("Tile should not be hazard!")
	}

	tile.SetHazard()
	if tile.IsHazard() != true {
		panic("Tile should have hazard!")
	}

	tile.ClearHazard()
	if tile.IsHazard() != false {
		panic("Hazard should be cleared!")
	}
}

func TestClear(t *testing.T) {
	tile := CreateEmptyTile()
	tile.SetHazard()
	tile.SetFood()

	if tile.IsHazard() && tile.IsFood() != true {
		panic("Tile Should have food and hazard!")
	}

	tile.Clear()
	if tile.IsFood() != false {
		panic("Tile should not have food!")
	}

	if tile.IsHazard() != true {
		panic("Tile hazard should not have been cleared")
	}
}

func TestStacked(t *testing.T) {
	tile := CreateDoubleStackTile(0, 0)

	if tile.IsDoubleStack() != true {
		panic("Should be a double stack!")
	}

	tile = CreateTripleStackTile(0)
	if tile.IsTripleStack() != true {
		panic("Should be a triple stack!")
	}

	if tile.IsSnakeSegment() != true {
		panic("Should be snake body part!")
	}

}
