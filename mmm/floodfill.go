package mmm

import (
	"sort"

	b "github.com/nosnaws/tiam/board"
)

func floodfill(board *b.FastBoard, startIndex, depth int, search []uint16) (int, int, int) {
	queue := []int{}
	visited := make(map[int]bool)

	queue = append(queue, startIndex)
	var current int
	depthMarker := -1
	currentDepth := 0
	foodDepth := -1
	searchDepth := -1
	queue = append(queue, depthMarker)

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current == depthMarker {
			currentDepth += 1
			queue = append(queue, depthMarker)

			if currentDepth >= depth || queue[0] == depthMarker {
				break
			}

		} else {
			if !visited[current] {
				visited[current] = true

				if board.IsTileFood(uint16(current)) && foodDepth == -1 {
					foodDepth = currentDepth
				}

				if isInSearch(search, uint16(current)) && searchDepth == -1 {
					searchDepth = currentDepth
				}

				index := uint16(current)
				for _, n := range board.GetNeighbors(index) {
					i := b.IndexInDirection(n, index, board.Width, board.Height, board.IsWrapped)
					queue = append(queue, int(i))
				}
			}
		}
	}

	total := 0
	for v := range visited {
		index := uint16(v)

		territoryBonus := len(board.GetNeighbors(index))
		total += 1 + territoryBonus
	}

	return total, foodDepth, searchDepth
}

func floodfillSafe(board *b.FastBoard, startIndex, depth int, minMoves []uint16) int {
	queue := []int{}
	visited := make(map[int]bool)

	queue = append(queue, startIndex)
	var current int
	depthMarker := -1
	currentDepth := 0
	queue = append(queue, depthMarker)

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current == depthMarker {
			currentDepth += 1
			queue = append(queue, depthMarker)

			if currentDepth >= depth || queue[0] == depthMarker {
				break
			}

		} else {
			if !visited[current] {
				visited[current] = true

				index := uint16(current)
				for _, n := range board.GetNeighbors(index) {
					i := b.IndexInDirection(n, index, board.Width, board.Height, board.IsWrapped)
					for _, m := range minMoves {
						if i == m {
							continue
						}
					}
					queue = append(queue, int(i))
				}
			}
		}
	}

	total := 0
	for v := range visited {
		index := uint16(v)

		territoryBonus := len(board.GetNeighbors(index))
		total += 1 + territoryBonus
	}

	return total
}

func findBF(board *b.FastBoard, startIndex int, testerFn func(int) bool) int {
	queue := []int{}
	visited := make(map[int]bool)

	queue = append(queue, startIndex)
	var current int
	depthMarker := -1
	currentDepth := 0
	searchDepth := -1
	queue = append(queue, depthMarker)

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if searchDepth != -1 {
			break
		}

		if current == depthMarker {
			currentDepth += 1
			queue = append(queue, depthMarker)

			if queue[0] == depthMarker {
				break
			}

		} else {
			if !visited[current] {
				visited[current] = true

				if testerFn(current) && searchDepth == -1 {
					searchDepth = currentDepth
				}

				index := uint16(current)
				for _, n := range board.GetNeighbors(index) {
					i := b.IndexInDirection(n, index, board.Width, board.Height, board.IsWrapped)
					queue = append(queue, int(i))
				}
			}
		}
	}

	return searchDepth
}

type DepthIndex struct {
	Index uint16
	Depth int
}

func findAllBF(board *b.FastBoard, startIndex int, testerFn func(int) bool) []DepthIndex {
	queue := []int{}
	visited := make(map[int]bool)

	queue = append(queue, startIndex)
	var current int
	depthMarker := -1
	currentDepth := 0
	searchResults := []DepthIndex{}
	queue = append(queue, depthMarker)

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current == depthMarker {
			currentDepth += 1
			queue = append(queue, depthMarker)

			if queue[0] == depthMarker {
				break
			}

		} else {
			if !visited[current] {
				visited[current] = true

				if testerFn(current) {

					searchResults = append(searchResults, DepthIndex{
						Index: uint16(current),
						Depth: currentDepth,
					})
				}

				index := uint16(current)
				for _, n := range board.GetNeighbors(index) {
					i := b.IndexInDirection(n, index, board.Width, board.Height, board.IsWrapped)
					queue = append(queue, int(i))
				}
			}
		}
	}

	return searchResults
}

func findBFUnsafe(board *b.FastBoard, startIndex int, testerFn func(int) bool) int {
	queue := []int{}
	visited := make(map[int]bool)

	queue = append(queue, startIndex)
	var current int
	depthMarker := -1
	currentDepth := 0
	searchDepth := -1
	queue = append(queue, depthMarker)

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if searchDepth != -1 {
			break
		}

		if current == depthMarker {
			currentDepth += 1
			queue = append(queue, depthMarker)

			if queue[0] == depthMarker {
				break
			}

		} else {
			if !visited[current] {
				visited[current] = true

				if testerFn(current) && searchDepth == -1 {
					searchDepth = currentDepth
				}

				index := uint16(current)
				for _, n := range board.GetNeighborsUnsafe(index) {
					i := b.IndexInDirection(n, index, board.Width, board.Height, board.IsWrapped)
					queue = append(queue, int(i))
				}
			}
		}
	}

	return searchDepth
}

func isInSearch(search []uint16, test uint16) bool {
	for _, s := range search {
		if s == test {
			return true
		}
	}
	return false
}

func dfSpaceFill(board *b.FastBoard, startIndex int, depth int) int {
	visited := make(map[int]bool)

	return dfSpaceFillRec(board, visited, startIndex, depth)
}

func dfSpaceFillRec(board *b.FastBoard, visited map[int]bool, node int, depth int) int {
	visited[node] = true

	total := 0

	if depth == 0 {
		return 1
	}

	index := uint16(node)
	neighbors := board.GetNeighbors(index)
	neighborI := []uint16{}
	for _, n := range neighbors {
		neighborI = append(neighborI, b.IndexInDirection(n, index, board.Width, board.Height, board.IsWrapped))
	}
	sort.Slice(neighborI, func(i, j int) bool {
		return len(board.GetNeighbors(neighborI[i])) < len(board.GetNeighbors(neighborI[j]))
	})

	for _, n := range neighborI {
		if !visited[int(n)] {
			t := dfSpaceFillRec(board, visited, int(n), depth-1)
			if t > total {
				total = t
			}
		}
	}

	return total
}
