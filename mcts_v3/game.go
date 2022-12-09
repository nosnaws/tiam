package mctsv3

import (
	api "github.com/nosnaws/tiam/battlesnake"
	"github.com/nosnaws/tiam/bitboard"
)

type SnakeMeta struct {
	Id   string
	Name string
}

type Game struct {
	Id     string
	Snakes []*SnakeMeta
	Me     *SnakeMeta
	Root   *node
}

type GameController struct {
	games map[string]*Game
}

func CreateGameController() *GameController {
	return &GameController{
		games: make(map[string]*Game),
	}
}

func (gc *GameController) StartGame(state api.GameState) {
	snakes := []*SnakeMeta{}
	for _, s := range state.Board.Snakes {
		snakes = append(snakes, createSnake(s))
	}

	gc.games[state.Game.ID] = &Game{
		Id:     state.Game.ID,
		Snakes: snakes,
		Me:     createSnake(state.You),
	}
}

func (gc *GameController) GetNextMove(state api.GameState, config MCTSConfig) bitboard.Dir {
	//game := gc.games[state.Game.ID]
	//newRoot := ChooseNewRoot(game.Root, state)
	//game.Root = newRoot

	//return MCTS(newRoot, config.MinimumSims, config.SimLimit, config.TreeDepthLimit)
	return GetNextMove(bitboard.CreateBitBoard(state), config, state.You.ID)
}

func (gc *GameController) EndGame(state api.GameState) {
	delete(gc.games, state.Game.ID)
}

func createSnake(snake api.Battlesnake) *SnakeMeta {
	return &SnakeMeta{
		Id:   snake.ID,
		Name: snake.Name,
	}
}
