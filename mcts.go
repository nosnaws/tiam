package main

import (
	"fmt"

	"github.com/BattlesnakeOfficial/rules"
)

type Node struct {
	youId              string
	ruleset            rules.Ruleset
	board              rules.BoardState
	scores             map[string]int32
	unexpandedChildren []*Node
	children           []*Node
	parent             *Node
}
