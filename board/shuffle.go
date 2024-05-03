package board

import (
	"math/rand"
)

func Shuffle[T any](arr []T) []T {
	arrToShuffle := make([]T, len(arr))
	copy(arrToShuffle, arr)

	rand.Shuffle(len(arrToShuffle), func(i, j int) {
		arrToShuffle[i], arrToShuffle[j] = arrToShuffle[j], arrToShuffle[i]
	})
	return arrToShuffle
}
