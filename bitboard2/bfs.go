package bitboard

import "github.com/nosnaws/tiam/moveset"

// returns the length of the path to the goal
// returns -1 if not found
func BFS(board *BitBoard, startIndex int, testerFn func(int) bool) int {
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

				index := current
				ms := board.GetNeighbors(index)
				if moveset.HasDown(ms) {
					i := indexInDirection(Down, index, board.width, board.height, board.isWrapped)
					queue = append(queue, int(i))

				}
				if moveset.HasUp(ms) {
					i := indexInDirection(Up, index, board.width, board.height, board.isWrapped)
					queue = append(queue, int(i))

				}
				if moveset.HasLeft(ms) {
					i := indexInDirection(Left, index, board.width, board.height, board.isWrapped)
					queue = append(queue, int(i))

				}
				if moveset.HasRight(ms) {
					i := indexInDirection(Right, index, board.width, board.height, board.isWrapped)
					queue = append(queue, int(i))

				}

			}
		}
	}

	return searchDepth
}
