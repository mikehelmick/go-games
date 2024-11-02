package controller

import "github.com/mikehelmick/go-games/pkg/game"

type GameController interface {
	CreateGame([]string) (string, error)
	LoadGame(string) (game.Game, error)
	Type() string
}
