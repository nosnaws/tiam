package moveset

import "math/rand"

// Uses 4 least significant bits for directions
// LEFT DOWN UP RIGHT
// 1 means the move is available
// 0 means it is not available

type MoveSet uint8
type Dir string

const (
	Left  Dir = "left"
	Right     = "right"
	Up        = "up"
	Down      = "down"
)

const leftMask = 0b00001000
const rightMask = 0b00000001
const downMask = 0b00000100
const upMask = 0b00000010

func Create() MoveSet {
	return 0b00000000
}

func SetLeft(ms MoveSet) MoveSet {
	return ms | leftMask
}

func SetRight(ms MoveSet) MoveSet {
	return ms | rightMask
}

func SetDown(ms MoveSet) MoveSet {
	return ms | downMask
}

func SetUp(ms MoveSet) MoveSet {
	return ms | upMask
}

func HasUp(ms MoveSet) bool {
	return ms&upMask > 0
}

func HasDown(ms MoveSet) bool {
	return ms&downMask > 0
}

func HasLeft(ms MoveSet) bool {
	return ms&leftMask > 0
}

func HasRight(ms MoveSet) bool {
	return ms&rightMask > 0
}

func IsEmpty(ms MoveSet) bool {
	return ms == 0
}

// possible bug where you pass a blank moveset in and get true back
func CountMoves(ms MoveSet) int {
	moves := 0

	if HasLeft(ms) {
		moves += 1
	}
	if HasRight(ms) {
		moves += 1
	}
	if HasUp(ms) {
		moves += 1
	}
	if HasDown(ms) {
		moves += 1
	}

	return moves
}

func Split(ms MoveSet) []MoveSet {
	moves := []MoveSet{}

	if HasLeft(ms) {
		moves = append(moves, SetLeft(Create()))
	}
	if HasRight(ms) {
		moves = append(moves, SetRight(Create()))
	}
	if HasUp(ms) {
		moves = append(moves, SetUp(Create()))
	}
	if HasDown(ms) {
		moves = append(moves, SetDown(Create()))
	}

	return moves
}

func ToDirs(ms MoveSet) []Dir {
	dirs := []Dir{}

	if HasLeft(ms) {
		dirs = append(dirs, Left)
	}
	if HasRight(ms) {
		dirs = append(dirs, Right)
	}
	if HasDown(ms) {
		dirs = append(dirs, Down)
	}
	if HasUp(ms) {
		dirs = append(dirs, Up)
	}

	return dirs
}

func GetRandomMove(ms MoveSet, rand *rand.Rand) MoveSet {
	numMoves := CountMoves(ms)
	if numMoves < 2 {
		return ms
	}

	for {
		n := rand.Intn(3)

		if n == 0 && HasLeft(ms) {
			return SetLeft(Create())
		}
		if n == 1 && HasDown(ms) {
			return SetDown(Create())
		}
		if n == 2 && HasUp(ms) {
			return SetUp(Create())
		}
		if n == 3 && HasRight(ms) {
			return SetRight(Create())
		}
	}
}
