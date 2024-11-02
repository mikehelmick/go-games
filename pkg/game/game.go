package game

type Game interface {
	String() string

	SetID(string)
	ID() string

	HandleInput(id string, command string) error

	Type() string

	// Sets the state of the players from the database
	SetPlayers([]string)
	// Gets the state of the players that should be saved in the database
	GetPlayers() []string

	// State returns the state to either load into or save from.
	State() interface{}

	// Winner should return the ID of the player who has won.
	Winner() string
}
