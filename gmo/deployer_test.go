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
	//headers := []string{
	//"EXPLORATION_CONSTANT=1.8",
	//"ALPHA_CONSTANT=0.531146",
	//"VORONOI_CONSTANT=1.000000",
	//"FOOD_CONSTANT_A=1.827048",
	//"FOOD_CONSTANT_B=14.860363",
	//"BIG_SNAKE_CONSTANT=16.095620",
	//}

	d := createDeployment(9090)
	//d.addContainer("tiam-test", "tiam", []string{})
	d.addContainer("tiam-main", "tiam-main", []string{})
	//d.addContainer("mini-test", "mini", []string{})
	d.addContainer("eater-test", "eater", []string{})

	ctx := context.Background()
	isSuccess := d.run(ctx)
	if !isSuccess {
		panic("Did not deploy successfully")
	}
	defer d.stopAndRemoveContainers(ctx)

	g := game{
		mapName: "standard",
		width:   11,
		height:  11,
		cliExe:  "../../rules/battlesnake",
		players: []player{
			//{
			//name: "tiam-test",
			//url:  "http://localhost:8090",
			//},
			{
				name: "taim-main",
				url:  "http://localhost:9090",
			},
			{
				name: "eater-test",
				url:  "http://localhost:9091",
			},
			{
				name: "mini-test",
				url:  "http://localhost:8080",
			},
		},
	}

	scores := make(map[string]int)

	for i := 0; i < 1; i++ {
		out := runGame(g, true)
		fmt.Println(out)

		if out != "" {
			scores[out] += 1
		}

		fmt.Println("FINISHED GAME")
	}

	fmt.Println(scores)

}
