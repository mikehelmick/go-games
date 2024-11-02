package fakedatabase

import (
	"time"

	"github.com/mikehelmick/go-games/pkg/model"
)

type storage struct {
	Players   map[string]model.Player    `json:"players"`
	Games     map[string]model.GameState `json:"games"`
	UpdatedAt time.Time                  `json:"updatedAt"`
}

func (s *storage) init() {
	if s.Players == nil {
		s.Players = make(map[string]model.Player)
	}
	if s.Games == nil {
		s.Games = make(map[string]model.GameState)
	}
}
