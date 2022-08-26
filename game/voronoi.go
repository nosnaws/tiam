package game

type pair struct {
	id    SnakeId
	index uint16
}

type voronoiResult struct {
	Score     int8
	FoodDepth int
}

func Voronoi(game *FastBoard, player SnakeId) voronoiResult {

	// create queue of indices
	queue := []pair{}
	// create map of visited indices
	visited := make(map[uint16]SnakeId)
	// create map of scores for indices
	scores := make(map[SnakeId]int8)
	depth := 0
	foodDepth := -1

	depthMark := pair{id: SnakeId(0), index: 0}
	mark := SnakeId(0)

	//
	// add snake heads to queue
	// mark snake heads as visited
	for id, head := range game.Heads {
		if game.IsSnakeAlive(id) {
			p := pair{id, head}
			queue = append(queue, p)
			visited[head] = id
		}
	}
	// add depth mark to queue
	queue = append(queue, depthMark)

	//
	// while the queue is not empty
	for len(queue) > 0 {
		//fmt.Println(queue)
		//fmt.Println(scores)
		//for y := int(game.height - 1); y >= 0; y-- {
		//var line string
		//for x := 0; x < int(game.width); x++ {
		//p := Point{X: int8(x), Y: int8(y)}
		//index := pointToIndex(p, game.width)
		//part := ""
		//if id, ok := visited[index]; ok {
		//part = fmt.Sprintf(" _%d ", id)
		//} else {
		//part = game.tileToString(index)
		//}
		//line = line + part
		//}
		//fmt.Println(line)
		//}
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
			neighbors := game.GetNeighbors(current.index)
			for _, neighbor := range neighbors {
				nIndex := IndexInDirection(neighbor, current.index, game.width, game.height, game.isWrapped)

				if game.IsTileFood(nIndex) && current.id == player && foodDepth == -1 {
					foodDepth = depth
				}
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

	return voronoiResult{Score: scores[player], FoodDepth: foodDepth}
}
