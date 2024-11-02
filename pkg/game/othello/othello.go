package othello

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/mikehelmick/go-games/pkg/game"
)

var _ game.Game = (*Game)(nil)

const TYPE = "othello"

type cell int

const (
	EMPTY cell = iota
	BLACK
	WHITE
)

type State struct {
	Board  [][]cell `json:"board"`
	Black  string   `json:"black"`
	White  string   `json:"white"`
	Winner string   `json:"winner,omitempty"`
}

type Game struct {
	id        string
	players   []string
	GameState State
}

var dirs []image.Point = []image.Point{
	{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1},
}

func initBoard() [][]cell {
	board := make([][]cell, 8)
	for i := 0; i < 8; i++ {
		board[i] = make([]cell, 8)
	}

	board[3][3] = WHITE
	board[3][4] = BLACK
	board[4][3] = BLACK
	board[4][4] = WHITE

	return board
}

func BlankGame() *Game {
	return &Game{
		id:      "",
		players: make([]string, 0, 2),
		GameState: State{
			Board:  initBoard(),
			Winner: "",
			Black:  "",
			White:  "",
		},
	}
}

func New(players []string) (*Game, error) {
	if len(players) != 2 {
		return nil, fmt.Errorf("game requires exactly 2 players")
	}
	if players[0] == players[1] {
		return nil, fmt.Errorf("players cannot be the same")
	}

	return &Game{
		id:      "",
		players: players,
		GameState: State{
			Board:  initBoard(),
			Winner: "",
			Black:  players[0],
			White:  players[1],
		},
	}, nil
}

func getPiece(cell cell) string {
	switch cell {
	case BLACK:
		return "🔴"
	case WHITE:
		return "🔵"
	default:
		return "🔶"
	}
}

func (g *Game) String() string {
	builder := strings.Builder{}
	for _, row := range g.GameState.Board {
		for _, cell := range row {
			builder.WriteString(getPiece(cell))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// Does not fully validate the moves... lazy.
func (g *Game) HandleInput(playerId string, command string) error {
	if g.GameState.Winner != "" {
		return fmt.Errorf("game is over, no moves allowed")
	}

	if g.players[0] != playerId {
		return fmt.Errorf("it is not player %q's turn", playerId)
	}

	parts := strings.Split(command, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid command, expect 2 ints separated by comma")
	}
	row, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("row %q is invalid: %w", parts[0], err)
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("col %q is invalid: %w", parts[1], err)
	}

	if g.GameState.Board[row][col] != EMPTY {
		return fmt.Errorf("must select an empty square")
	}

	color := WHITE
	if playerId == g.GameState.Black {
		color = BLACK
	}

	if g.flipCells(image.Point{row, col}, color) == 0 {
		return fmt.Errorf("invalid move")
	}

	// we have made a valid move, yay.
	// rotate players
	g.GameState.Board[row][col] = color
	g.players = []string{g.players[1], g.players[0]}
	// check winner
	whiteCount := 0
	blackCount := 0
	for _, row := range g.GameState.Board {
		for _, col := range row {
			if col == WHITE {
				whiteCount++
			} else if col == BLACK {
				blackCount++
			}
		}
	}
	if whiteCount+blackCount == 64 {
		if whiteCount > blackCount {
			g.GameState.Winner = g.GameState.White
		} else if blackCount > whiteCount {
			g.GameState.Winner = g.GameState.Black
		} else {
			g.GameState.Winner = "tie"
		}
	}

	return nil
}

func (g *Game) flipCells(start image.Point, color cell) int {
	flipped := 0
	for _, dir := range dirs {
		flipped += g.fillLine(start, dir, color)
	}
	return flipped
}

func valid(p image.Point) bool {
	return p.X >= 0 && p.X <= 8 && p.Y >= 0 && p.Y <= 8
}

func (g *Game) fillLine(start image.Point, dir image.Point, color cell) int {
	possible := make([]image.Point, 0)
	doFill := false // assume we will fill
	flipCount := 0

	next := start.Add(dir)
	if !valid(next) {
		return flipCount
	}
	if val := g.GameState.Board[next.X][next.Y]; val == color || val == EMPTY {
		// nothing along this path
		return flipCount
	}

	for valid(next) {
		if val := g.GameState.Board[next.X][next.Y]; val != EMPTY && val != color {
			// this is a valid possibility
			possible = append(possible, next)
		} else if val != EMPTY && val == color {
			if len(possible) > 0 {
				doFill = true
			}
			break
		} else if val == EMPTY {
			break
		}
		next = next.Add(dir)
	}

	if doFill {
		for _, c := range possible {
			g.GameState.Board[c.X][c.Y] = color
		}
		return len(possible)
	}
	return 0
}

func (g *Game) SetID(id string) {
	g.id = id
}

func (g *Game) ID() string {
	return g.id
}

func (g *Game) Type() string {
	return TYPE
}

func (g *Game) SetPlayers(players []string) {
	g.players = players
}

func (g *Game) GetPlayers() []string {
	return g.players
}

func (g *Game) State() interface{} {
	return &g.GameState
}

func (g *Game) Winner() string {
	return g.GameState.Winner
}
