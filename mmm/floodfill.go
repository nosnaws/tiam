package mmm

import (
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

	return len(visited), foodDepth, searchDepth
}

func findBF(board *b.FastBoard, startIndex int, testerFn func(b.Tile) bool) int {
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

				if testerFn(board.List[current]) && searchDepth == -1 {
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

func isInSearch(search []uint16, test uint16) bool {
	for _, s := range search {
		if s == test {
			return true
		}
	}
	return false
}
