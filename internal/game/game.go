package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"github.com/slonegd-go/reversi/internal/player"
)

type Game struct {
	cells     []player.Color
	stepCellN int
	players   []player.Player
	log       func(string, ...interface{})
}

type Options struct {
	log func(string, ...interface{})
}

type Option func(*Options)

func WithLogger(log func(string, ...interface{})) Option {
	return func(opts *Options) {
		opts.log = log
	}
}

func New(p1, p2 player.Player, opts ...Option) *Game {
	cells := make([]player.Color, 64)
	cells[27] = player.Green
	cells[28] = player.Red
	cells[35] = player.Red
	cells[36] = player.Green

	p1.SetColor(player.Green)
	p2.SetColor(player.Red)

	options := &Options{
		log: func(string, ...interface{}) {},
	}

	for _, opt := range opts {
		opt(options)
	}

	game := &Game{
		cells:     cells,
		stepCellN: -1,
		players:   []player.Player{p1, p2},
		log:       options.log,
	}

	game.log(game.String())

	return game
}

func (game *Game) Start() string {
	game.log(game.String())
	for i := 0; i < 64; i++ {
		currentPlayer := game.players[i%2]

		game.log("%s player step:", currentPlayer.Color())
		currentPlayer.Step(game.cells, func(position string) error {
			err := game.Step(currentPlayer.Color(), position)
			if err != nil {
				game.log(err.Error())
			}
			return err
		})

		end := game.endCheck(currentPlayer.Color())
		if end {
			winPlayer, losePlayer := game.compute()
			result := fmt.Sprintf("%s player win", winPlayer.Color())
			game.log(result)
			winPlayer.Notify(player.Win)
			losePlayer.Notify(player.Lose)
			return result
		}
	}
	return "error"
}

func (game *Game) endCheck(color player.Color) bool {
	otherPlayer := player.Green
	if color == player.Green {
		otherPlayer = player.Red
	}
	return !game.hasSteps(otherPlayer)
}

func (game *Game) compute() (win player.Player, lose player.Player) {
	win = game.players[0]
	lose = game.players[1]
	greenCount, redCount := 0, 0
	for _, cell := range game.cells {
		switch cell {
		case player.Green:
			greenCount++
		case player.Red:
			redCount++
		}
	}
	game.log("%s %d:%d %s", green("green"), greenCount, redCount, red("red"))
	if redCount > greenCount {
		win, lose = lose, win
	}
	return win, lose
}

func (game *Game) hasSteps(color player.Color) bool {
	for i, cell := range game.cells {
		if cell != player.Empty {
			continue
		}
		for _, direction := range directionList {
			if game.count(i, direction, color) != 0 {
				return true
			}
		}
	}
	return false
}

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func (game Game) String() string {
	var builder strings.Builder
	builder.Grow(101)
	builder.WriteString("\n\\ A B C D E F G H\n")
	for i := 0; i < 8; i++ {
		builder.WriteString(strconv.Itoa(i + 1))
		for j := 0; j < 8; j++ {
			number := i*8 + j
			if game.stepCellN == number {
				builder.WriteString(" x")
				continue
			}
			switch game.cells[i*8+j] {
			case player.Empty:
				builder.WriteString("  ")
			case player.Green:
				builder.WriteString(green(" G"))
			case player.Red:
				builder.WriteString(red(" R"))
			}
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (game *Game) Step(color player.Color, position string) error {
	if color != player.Green && color != player.Red {
		return errors.New("only green and red state available")
	}

	cellN, err := parseCellN(position)
	if err != nil {
		return fmt.Errorf("parse cell number: %w", err)
	}

	game.stepCellN = cellN

	game.log(game.String())

	if game.cells[cellN] != player.Empty {
		return errors.New("cell not empty")
	}

	directions := []direction{}
	for _, direction := range directionList {
		if game.count(cellN, direction, color) != 0 {
			directions = append(directions, direction)
		}
	}
	if len(directions) == 0 {
		return errors.New("unavailable step")
	}

	game.stepCellN = -1
	game.cells[cellN] = color
	for _, direction := range directions {
		game.count(cellN, direction, color, game.changeTo(color))
	}

	game.log(game.String())
	return nil
}

func (game *Game) changeTo(color player.Color) func(i int) {
	return func(i int) {
		game.cells[i] = color
	}
}

func (game *Game) count(cellN int, direction direction, color player.Color, change ...func(i int)) int {
	changeFunc := func(_ int) {}
	if len(change) != 0 {
		changeFunc = change[0]
	}
	count := 0

	switch direction {
	case up:
		if cellN < 8 { // border
			return 0
		}
		i := cellN - 8
		for ; i < 0 || game.cells[i] != color; i -= 8 {
			if i < 8 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case down:
		if cellN > 63-8 { // border
			return 0
		}
		i := cellN + 8
		for ; i > 63 || game.cells[i] != color; i += 8 {
			if i > 63-8 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	case left:
		if cellN%8 == 0 { // border
			return 0
		}
		i := cellN - 1
		for ; i < 0 || game.cells[i] != color; i-- {
			if i < 0 || i%8 == 0 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case right:
		if cellN%8 == 7 { // border
			return 0
		}
		i := cellN + 1
		for ; i > 63 || game.cells[i] != color; i++ {
			if i > 63 || i%8 == 7 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	case leftup:
		if cellN < 8 || cellN%8 == 0 { // border
			return 0
		}
		i := cellN - 9
		for ; i < 0 || game.cells[i] != color; i -= 9 {
			if i%8 == 0 || i < 8 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case rightup:
		if cellN < 8 || cellN%8 == 7 { // border
			return 0
		}
		i := cellN - 7
		for ; i < 0 || game.cells[i] != color; i -= 7 {
			if i%8 == 7 || i < 8 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case rightdown:
		if cellN > 63-8 || cellN%8 == 7 { // border
			return 0
		}
		i := cellN + 9
		for ; i > 63 || game.cells[i] != color; i += 9 {
			if i%8 == 7 || i > 63-8 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	case leftdown:
		if cellN > 63-8 || cellN%8 == 0 { // border
			return 0
		}
		i := cellN + 7
		for ; i > 63 || game.cells[i] != color; i += 7 {
			if i%8 == 0 || i > 63-8 || game.cells[i] == player.Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	} // switch direction

	return count
}

func parseCellN(position string) (int, error) {
	if len(position) != 2 {
		return 0, fmt.Errorf("position must be from A1 to H8, got: %s", position)
	}
	column := position[0] - byte('A')
	line := position[1] - byte('1')
	if column > 7 || line > 7 {
		return 0, fmt.Errorf("position must be from A1 to H8, got: %s", position)
	}

	cellN := column + line*8
	return int(cellN), nil
}

type direction int

const (
	left direction = iota
	right
	up
	down
	leftup
	rightup
	leftdown
	rightdown
)

var directionList = []direction{left, right, up, down, leftup, rightup, leftdown, rightdown}
