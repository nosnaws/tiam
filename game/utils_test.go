package game

import (
	"fmt"
	"testing"
)

func TestIndexInDirection(t *testing.T) {

	if indexInDirection(Right, 0, 3) != 1 {
		panic("Right should be TileIndex 1")
	}

	if indexInDirection(Up, 0, 3) != 3 {
		panic("Up should be TileIndex 3")
	}

	if indexInDirection(Down, 4, 3) != 1 {
		panic("Down should be TileIndex 1")
	}

	if indexInDirection(Left, 1, 3) != 0 {
		fmt.Println(indexInDirection(Left, 1, 3))
		panic("Left should be TileIndex 0")
	}

}
