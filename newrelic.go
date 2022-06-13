package main

import (
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func getCustomAttributesEnd(txn *newrelic.Transaction, state GameState) {
	snakes := state.Board.Snakes

	txn.AddAttribute("snakeGameId", state.Game.ID)
	txn.AddAttribute("snakeRules", state.Game.Ruleset.Name)
	txn.AddAttribute("snakeTurn", state.Turn)

	txn.AddAttribute("snakeName", state.You.Name)
	txn.AddAttribute("snakeId", state.You.ID)
	txn.AddAttribute("snakeHealth", state.You.Health)
	txn.AddAttribute("snakeLength", state.You.Length)

	var winnerName string
	var winnerId string
	var isWinner bool
	if len(snakes) > 0 {
		winner := snakes[0]
		winnerName = winner.Name
		winnerId = winner.ID
		isWinner = winner.Name == state.You.Name
	}

	replayLink := fmt.Sprintf("https://play.battlesnake.com/g/%s", state.Game.ID)
	txn.AddAttribute("snakeGameWinnerName", winnerName)
	txn.AddAttribute("snakeGameWinnerId", winnerId)
	txn.AddAttribute("snakeGameIsWin", isWinner)
	txn.AddAttribute("snakeGameReplayLink", replayLink)
}
