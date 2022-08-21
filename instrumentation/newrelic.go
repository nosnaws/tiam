package instrumentation

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	fastGame "github.com/nosnaws/tiam/game"
)

type SnakeData struct {
	Head  fastGame.Coord   `json:"head"`
	Body  []fastGame.Coord `json:"body"`
	Color string           `json:"color"`
}

func encodeString(s []byte) string {
	return b64.StdEncoding.EncodeToString(s)
}

func getSnakeData(snake fastGame.Battlesnake) SnakeData {
	return SnakeData{Head: snake.Head, Body: snake.Body, Color: snake.Customizations.Color}
}

func getBaseAttributes(txn *newrelic.Transaction, state fastGame.GameState) {
	snakeData, _ := json.Marshal(getSnakeData(state.You))
	foodString, _ := json.Marshal(state.Board.Food)
	hazardString, _ := json.Marshal(state.Board.Hazards)

	txn.AddAttribute("snakeGameId", state.Game.ID)
	txn.AddAttribute("snakeRules", state.Game.Ruleset.Name)
	txn.AddAttribute("snakeTurn", state.Turn)
	txn.AddAttribute("snakeMap", state.Game.Map)

	txn.AddAttribute("snakeName", state.You.Name)
	txn.AddAttribute("snakeId", state.You.ID)
	txn.AddAttribute("snakeHealth", state.You.Health)
	txn.AddAttribute("snakeLength", state.You.Length)
	txn.AddAttribute("snakeData", encodeString(snakeData))

	txn.AddAttribute("snakeBoardHeight", state.Board.Height)
	txn.AddAttribute("snakeBoardWidth", state.Board.Width)
	txn.AddAttribute("snakeBoardFood", encodeString(foodString))
	txn.AddAttribute("snakeBoardHazards", encodeString(hazardString))

	for i, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		opponentSnakeData, _ := json.Marshal(getSnakeData(snake))

		snakeIndex := i + 1
		txn.AddAttribute(fmt.Sprintf("snakeOpponent_%d_Name", snakeIndex), snake.Name)
		txn.AddAttribute(fmt.Sprintf("snakeOpponent_%d_Id", snakeIndex), snake.ID)
		txn.AddAttribute(fmt.Sprintf("snakeOpponent_%d_Health", snakeIndex), snake.Health)
		txn.AddAttribute(fmt.Sprintf("snakeOpponent_%d_Length", snakeIndex), snake.Length)
		txn.AddAttribute(fmt.Sprintf("snakeOpponent_%d_Data", snakeIndex), encodeString(opponentSnakeData))
	}
}

func GetCustomAttributes(txn *newrelic.Transaction, state fastGame.GameState) {
	getBaseAttributes(txn, state)
}

func GetCustomAttributesEnd(txn *newrelic.Transaction, state fastGame.GameState) {
	getBaseAttributes(txn, state)
	snakes := state.Board.Snakes

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
