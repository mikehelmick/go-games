package fakedatabase

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mikehelmick/go-games/pkg/game"
	"github.com/mikehelmick/go-games/pkg/model"
)

const file = "database.json"

type Database struct {
	data storage
	mu   sync.RWMutex
}

func Init() *Database {
	db := Database{}
	db.data.init()

	bytes, err := os.ReadFile(file)
	if err == nil {
		if err := json.Unmarshal(bytes, &db.data); err != nil {
			log.Fatalf("database corrupted... sorry: %v", err)
		}
	} else {
		log.Print("no database found, initializing")
	}

	return &db
}

// save must be called under a write lock.
func (db *Database) save() error {
	db.data.UpdatedAt = time.Now().UTC()
	bytes, err := json.MarshalIndent(&db.data, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal data: %w", err)
	}
	if err := os.WriteFile(file, bytes, 0666); err != nil {
		log.Fatalf("error writing database: %v", err)
	}
	log.Printf("wrote database at %v", db.data.UpdatedAt.Format(time.RFC3339))
	return nil
}

func (db *Database) GetPlayer(id string) (*model.Player, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	player, ok := db.data.Players[id]
	if ok {
		return player.Clone(), nil
	}
	return nil, fmt.Errorf("player not found")
}

func (db *Database) FindPlayerByEmail(email string) (*model.Player, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, v := range db.data.Players {
		if v.Email == email {
			return v.Clone(), nil
		}
	}
	return nil, fmt.Errorf("player not found")
}

func (db *Database) AddPlayer(p *model.Player) (*model.Player, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, v := range db.data.Players {
		if v.Email == p.Email {
			return v.Clone(), fmt.Errorf("user already exists")
		}
	}

	id := uuid.New().String()
	p.ID = id
	db.data.Players[id] = *p

	if err := db.save(); err != nil {
		return nil, fmt.Errorf("unable to add player: %w", err)
	}
	return p, nil
}

func (db *Database) SaveGame(game game.Game) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if game.ID() == "" {
		game.SetID(uuid.NewString())
	}

	bytes, err := json.Marshal(game.State())
	if err != nil {
		return fmt.Errorf("cannot marshal game state: %w", err)
	}

	gameState := model.GameState{
		ID:       game.ID(),
		GameType: game.Type(),
		Players:  game.GetPlayers(),
		State:    string(bytes),
	}

	if old, ok := db.data.Games[game.ID()]; ok {
		if old.GameType != game.Type() {
			return fmt.Errorf("incorrect game type")
		}
	}
	db.data.Games[gameState.ID] = gameState

	return db.save()
}

func (db *Database) GetGame(id string, game game.Game) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, ok := db.data.Games[id]
	if !ok {
		return fmt.Errorf("invalid game ID")
	}
	if data.GameType != game.Type() {
		return fmt.Errorf("game id %q is of type %q which doesn't match game loader of type %q",
			id, data.GameType, game.Type())
	}

	// load the data into the game - defensive copies
	players := make([]string, len(data.Players))
	copy(players, data.Players)
	game.SetPlayers(players)
	game.SetID(id)

	// deserialize state into the game
	if err := json.Unmarshal([]byte(data.State), game.State()); err != nil {
		return fmt.Errorf("unable to unmarshal data into game: %w", err)
	}
	return nil
}
