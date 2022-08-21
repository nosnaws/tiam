package gmo

import "fmt"

func run() {
	g := game{
		mapName: "standard",
		width:   11,
		height:  11,
		cliExe:  "../rules/battlesnake",
		players: []player{
			{
				name: "test1",
				url:  "localhost:8080",
			},
			{
				name: "test2",
				url:  "localhost:8081",
			},
		},
	}

	out := runGame(g, false)

	fmt.Println(out)
}
