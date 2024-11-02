package model

type GameState struct {
	// ID of this game
	ID string `json:"id"`
	// GameType is which ID
	GameType string `json:"gameType"`
	// ID of the players, in turn order
	Players []string `json:"players"`
	// Serialized state of the game.
	State string `json:"state"`
}
