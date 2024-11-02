package main

import (
	"log"

	"github.com/mikehelmick/go-games/pkg/controller/othello"
	"github.com/mikehelmick/go-games/pkg/fakedatabase"
	"github.com/mikehelmick/go-games/pkg/model"
	"github.com/mikehelmick/go-games/pkg/server"
)

func main() {
	db := fakedatabase.Init()
	svr := server.GetServer()
	svr.ProvideDatabase(db)
	// could register more games here.
	svr.MustRegisterGame(othello.New(db))

	player1, player2 := createPlayers(db)

	// create a game
	id, err := svr.CreateGame([]string{player1, player2}, "othello")
	if err != nil {
		log.Fatalf("cannot greate game: %v", err)
	}

	// take some moves, would be like input coming from a UI layer.
	svr.UserInput(player1, "othello", id, "3,2")
	svr.UserInput(player2, "othello", id, "2,4")
	svr.UserInput(player1, "othello", id, "3,5")
	svr.UserInput(player2, "othello", id, "2,4")
}

func createPlayers(db *fakedatabase.Database) (string, string) {
	data := []*model.Player{
		{Name: "bob", Email: "bob@example.com"},
		{Name: "alice", Email: "alice@example.com"},
	}

	// create or load players
	for i := range data {
		p, err := db.FindPlayerByEmail(data[i].Email)
		if p != nil && err == nil {
			data[i] = p
			continue
		}

		updated, err := db.AddPlayer(data[i])
		if err != nil {
			panic(err)
		}
		data[i] = updated
	}

	return data[0].ID, data[1].ID
}
