package bitboard

import "github.com/nosnaws/tiam/moveset"

type pair struct {
	id    string
	index int
}

type VoronoiResult struct {
	Score     map[string]int
	FoodDepth map[string]int
	Territory map[uint16]string
}

func (bb *BitBoard) Voronoi() VoronoiResult {

	// create queue of indices
	queue := []pair{}
	// create map of visited indices
	visited := make(map[int]string)
	// create map of scores for indices
	scores := make(map[string]int)
	depth := 0
	//food := make(map[string]int)

	depthMark := pair{id: "depth mark", index: 0}
	mark := "mark"

	//
	// add snake heads to queue
	// mark snake heads as visited
	for id, snake := range bb.Snakes {
		//food[id] = -1
		if bb.IsSnakeAlive(id) {
			p := pair{id, snake.GetHeadIndex()}
			queue = append(queue, p)
			visited[p.index] = id
		}
	}
	// add depth mark to queue
	queue = append(queue, depthMark)

	//
	// while the queue is not empty
	for len(queue) > 0 {
		//    dequeue index
		var current pair
		current, queue = queue[0], queue[1:]
		//
		//    if index is depth mark
		if current == depthMark {
			//      increase depth count by 1
			depth += 1
			//      add depth mark to queue
			queue = append(queue, depthMark)
			//
			//      if front of queue is a depth mark
			if queue[0] == depthMark {
				//        end (we have searched all tiles)
				//
				break
			}
		} else { //    else
			// if this index is already marked, skip out
			if id, ok := visited[current.index]; ok && id == mark {
				continue
			}

			//      loop through neighbors of index
			neighbors := moveset.Split((bb.GetNeighbors(current.index)))
			for _, neighbor := range neighbors {
				nIndex := indexInDirection(MoveSetToDir(neighbor), current.index, bb.width, bb.height, bb.isWrapped)

				//if bb.IsTileFood(nIndex) && food[current.id] == -1 {
				//food[current.id] = depth
				//}
				//        if neighbor is in visited map
				if other, ok := visited[nIndex]; ok {
					//          if visited map is not a mark and the visited snake does not equal the snake for the current index
					if other != mark && other != current.id {
						//            reduce the score for this neighbor by one
						scores[other] -= 1
						//            set the visited map to mark for this neighbor
						visited[nIndex] = mark
					}
					//
				} else {
					//          increase the score for this neighbor by 1
					//territoryBonus := len(game.GetNeighbors(nIndex))
					scores[current.id] += 1
					//          add neighbor to visited map
					p := pair{id: current.id, index: nIndex}
					visited[nIndex] = current.id
					//          add neighbor to the queue
					queue = append(queue, p)
				}
			}
		}
	}

	return VoronoiResult{Score: scores}
}
