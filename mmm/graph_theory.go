package mmm

import (
	"github.com/nosnaws/tiam/board"
)

func getArticulationPoints(b *board.FastBoard, i uint16) []uint16 {
	parent := make(map[uint16]uint16)
	visited := make(map[uint16]bool)
	depth := make(map[uint16]int)
	low := make(map[uint16]int)

	return articulationHelper(b, parent, visited, depth, low, i, 0)
}

func minInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func articulationHelper(b *board.FastBoard, parent map[uint16]uint16, visited map[uint16]bool, depth, low map[uint16]int, i uint16, d int) []uint16 {
	visited[i] = true
	depth[i] = d
	low[i] = d
	childCount := 0
	isArticulation := false
	p := []uint16{}

	for _, ni := range b.GetNeighborIndices(i) {
		if _, ok := visited[ni]; !ok {
			parent[ni] = i
			newP := articulationHelper(b, parent, visited, depth, low, ni, d+1)
			for _, np := range newP {
				p = append(p, np)
			}

			childCount += 1

			if low[ni] >= depth[i] {
				isArticulation = true
			}
			low[i] = minInt(low[i], depth[ni])
		} else if ni != parent[i] {
			low[i] = minInt(low[i], depth[ni])
		}
	}

	_, pOk := parent[i]

	if (pOk && isArticulation) || (!pOk && childCount > 1) {
		p = append(p, i)
	}

	return p
}
