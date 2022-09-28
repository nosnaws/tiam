package gmo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//../rules/battlesnake play -W 11 -H 11 -n tiam -u http://localhost:8081 -n local -u http://localhost:8080 -n shai -u http://localhost:8082 -n shai2 -u http://localhost:8083 --output game_out.txt

type game struct {
	players                   []player
	width, height, hzDamage   int
	mapName, gameType, cliExe string
}

type player struct {
	name, url string
}

func runGame(g game, visual bool) string {
	cliExe, isExeSet := os.LookupEnv("BS_EXE")
	if !isExeSet {
		log.Fatal("Environment variable BS_EXE not set, exiting...")
	}

	commandPath, err := exec.LookPath(cliExe)
	if err != nil {
		log.Panic(err)
	}

	playerOptions := createPlayerOptions(g.players)
	playOptions := []string{
		"play",
		createMapHeightOption(g.height),
		createMapWidthOption(g.width),
		createMapOption(g.mapName),
		createGameTypeOption(g.gameType),
		createHazardDamageOption(g.hzDamage),
		"-t=500",
	}

	for _, p := range playerOptions {
		playOptions = append(playOptions, p)
	}

	if visual {
		playOptions = append(playOptions, "--browser")
	}

	command := exec.Command(commandPath, playOptions...)

	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR", string(output))
		log.Panic(err)
	}

	return parseWinner(string(output))
}

func parseWinner(out string) string {
	winnerRegex := regexp.MustCompile(`Game completed after .+ turns. (.+) was the winner.`)
	isDraw := strings.Contains(out, "It was a draw.")

	if isDraw {
		return ""
	}

	winner := winnerRegex.FindStringSubmatch(out)
	fmt.Println(winner)

	return winner[1]
}

func createPlayerOptions(players []player) []string {
	var options []string
	for _, p := range players {
		opt := []string{"-n=" + p.name, "-u=" + p.url}
		options = append(options, opt...)
	}

	return options
}

func createMapOption(name string) string {
	return fmt.Sprintf("-m=%s", name)
}

func createMapWidthOption(width int) string {
	return fmt.Sprintf("-W=%d", width)
}

func createMapHeightOption(height int) string {
	return fmt.Sprintf("-H=%d", height)
}

func createGameTypeOption(name string) string {
	return fmt.Sprintf("-g=%s", name)
}

func createHazardDamageOption(damage int) string {
	return fmt.Sprintf("--hazardDamagePerTurn=%d", damage)
}
