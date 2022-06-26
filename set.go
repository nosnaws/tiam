package main

import "github.com/BattlesnakeOfficial/rules"

type Set struct {
	list map[rules.Point]struct{}
}

func CreateSet(points []rules.Point) Set {
	list := make(map[rules.Point]struct{})

	s := Set{list: list}
	for _, p := range points {
		s.Add(p)
	}

	return s
}

func (x Set) Has(p rules.Point) bool {
	_, ok := x.list[p]
	return ok
}

func (x Set) Add(p rules.Point) {
	if !x.Has(p) {
		x.list[p] = struct{}{}
	}
}
