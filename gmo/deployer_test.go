package gmo

import (
	"context"
	"fmt"
	"testing"
)

func TestDeployDockerContainers(t *testing.T) {
	//t.Skip()

	//exploration := fmt.Sprintf("EXPLORATION_CONSTANT=%f", g[0])
	//alpha := fmt.Sprintf("ALPHA_CONSTANT=%f", g[1])
	//voronoi := fmt.Sprintf("VORONOI_CONSTANT=%f", g[2])
	//foodA := fmt.Sprintf("FOOD_CONSTANT_A=%f", g[3])
	//foodB := fmt.Sprintf("FOOD_CONSTANT_B=%f", g[3])
	//biggest := fmt.Sprintf("BIG_SNAKE_CONSTANT=%f", g[3])

	// specially-devoted-javelin
	headers := []string{
		//"EXPLORATION_CONSTANT=1.8",
		//"ALPHA_CONSTANT=0.531146",
		//"VORONOI_CONSTANT=1.000000",
		//"FOOD_CONSTANT_A=1.827048",
		//"FOOD_CONSTANT_B=14.860363",
		//"BIG_SNAKE_CONSTANT=16.095620",
		"GOMAXPROCS=1",
	}

	d := createDeployment(9091)
	//d.addContainer("tiam-test", "tiam", []string{})
	//d.addContainer("mini-test", "mini", []string{})
	d.addContainer("tiam-main", "tiam-main", headers)
	d.addContainer("eater-test", "eater", headers)
	d.addContainer("eater-test2", "eater", headers)

	ctx := context.Background()
	isSuccess := d.run(ctx)
	if !isSuccess {
		panic("Did not deploy successfully")
	}
	defer d.stopAndRemoveContainers(ctx)

	g := game{
		mapName:  "hz_islands_bridges",
		gameType: "wrapped",
		hzDamage: 100,
		width:    11,
		height:   11,
		cliExe:   "../../rules/battlesnake",
		players: []player{
			//{
			//name: "tiam-test",
			//url:  "http://localhost:9090",
			//},
			{
				name: "taim-main",
				url:  "http://localhost:9091",
			},
			{
				name: "eater-test",
				url:  "http://localhost:9092",
			},
			{
				name: "eater-test2",
				url:  "http://localhost:9092",
			},
			{
				name: "tiam-test",
				url:  "http://localhost:8080",
			},
		},
	}

	scores := make(map[string]int)

	inBrowser := true
	numGames := 1
	for i := 0; i < numGames; i++ {
		fmt.Println("Running game #", i)
		out := runGame(g, inBrowser)
		fmt.Println(out)

		if out != "" {
			scores[out] += 1
		}

		fmt.Println("FINISHED GAME")
	}

	fmt.Println(scores)

}
