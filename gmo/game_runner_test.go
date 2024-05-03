package gmo

import (
	"fmt"
	"testing"
)

func TestCreateGameRunCommand(t *testing.T) {
	t.Skip()
	g := game{
		mapName: "standard",
		width:   11,
		height:  11,
		cliExe:  "../../rules/battlesnake",
		players: []player{
			{
				name: "eater",
				url:  "http://localhost:8082",
			},
			{
				name: "random",
				url:  "http://localhost:8081",
			},
		},
	}

	out := runGame(g, true)

	fmt.Println(out)
}
