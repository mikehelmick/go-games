package server

import (
	"fmt"
	"log"
	"sync"

	"github.com/mikehelmick/go-games/pkg/controller"
	"github.com/mikehelmick/go-games/pkg/fakedatabase"
	"github.com/mikehelmick/go-games/pkg/model"
)

type Server struct {
	controllers map[string]controller.GameController
	db          *fakedatabase.Database
}

var (
	server *Server = &Server{
		controllers: make(map[string]controller.GameController),
	}
	mu sync.Mutex
)

func GetServer() *Server {
	return server
}

func (s *Server) ProvideDatabase(db *fakedatabase.Database) {
	s.db = db
}

func (s *Server) MustRegisterGame(gc controller.GameController) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := server.controllers[gc.Type()]; ok {
		panic(fmt.Sprintf("game controller %q already registered", gc.Type()))
	}
	server.controllers[gc.Type()] = gc
}

func (s *Server) RegisterUser(name string, email string) (string, error) {
	player, err := server.db.AddPlayer(&model.Player{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return "", err
	}
	return player.ID, nil
}

func (s *Server) CreateGame(players []string, typ string) (string, error) {
	gc := s.controllers[typ]
	if gc == nil {
		return "", fmt.Errorf("unknown game type: %q", typ)
	}

	id, err := gc.CreateGame(players)
	if err != nil {
		return "", fmt.Errorf("game create error: %w", err)
	}

	return id, nil
}

func (s *Server) UserInput(pID, gameType, gID, input string) {
	gc := s.controllers[gameType]
	if gc == nil {
		log.Printf("invalid game type: %q", gameType)
		return
	}

	game, err := gc.LoadGame(gID)
	if err != nil {
		log.Printf("unable to load game: %v", err)
		return
	}

	if err := game.HandleInput(pID, input); err != nil {
		log.Printf("%v", err)
		return
	}

	s.db.SaveGame(game)

	// game did it's thing.
	outout := game.String()
	log.Printf("\nplayer: %v took turn\ngame: %v\nstate\n%s\n", pID, gID, outout)
}
