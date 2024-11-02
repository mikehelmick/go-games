package othello

import (
	"log"

	"github.com/mikehelmick/go-games/pkg/controller"
	"github.com/mikehelmick/go-games/pkg/fakedatabase"
	"github.com/mikehelmick/go-games/pkg/game"
	"github.com/mikehelmick/go-games/pkg/game/othello"
)

var _ controller.GameController = (*OthelloController)(nil)

type OthelloController struct {
	database *fakedatabase.Database
}

func New(db *fakedatabase.Database) *OthelloController {
	return &OthelloController{
		database: db,
	}
}

func (c *OthelloController) Type() string {
	return othello.TYPE
}

func (c *OthelloController) CreateGame(players []string) (string, error) {
	newGame, err := othello.New(players)
	if err != nil {
		return "", err
	}

	if err := c.database.SaveGame(newGame); err != nil {
		return "", err
	}

	log.Printf("CREATED GAME: type: othello ID: %v\n%s\n", newGame.ID(), newGame.String())

	return newGame.ID(), nil
}

func (c *OthelloController) LoadGame(gameID string) (game.Game, error) {
	// load the game
	game := othello.BlankGame()
	if err := c.database.GetGame(gameID, game); err != nil {
		return nil, err
	}
	return game, nil
}
