package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"github.com/slonegd-go/reversi/internal/player"
)

type State int

const (
	Empty State = iota
	Green
	Red
)

func (state State) not(other State) bool {
	return state != Empty && state != other
}

type Game struct {
	cells     [64]State
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
	cells := [64]State{}
	cells[27] = Green
	cells[28] = Red
	cells[35] = Red
	cells[36] = Green

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

func (game *Game) Start() {
	game.log(game.String())
	for i := 0; i < 60; i++ {
		currentPlayer := game.players[i%2]
		for {
			game.log("%s player step:", currentPlayer.Color())
			position := currentPlayer.Step(nil)
			err := game.Step(State(currentPlayer.Color()), position)
			if err != nil {
				game.log(err.Error())
				continue
			}
			break
		}
		ok := game.endCheck(currentPlayer.Color())
		if ok {
			currentPlayer.Notify(player.Win)
			return
		}
	}
}

func (game *Game) endCheck(player.Color) bool {
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
			case Empty:
				builder.WriteString("  ")
			case Green:
				builder.WriteString(green(" G"))
			case Red:
				builder.WriteString(red(" R"))
			}
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (game *Game) Step(color State, position string) error {
	if color != Green && color != Red {
		return errors.New("only green and red state available")
	}

	cellN, err := parseCellN(position)
	if err != nil {
		return fmt.Errorf("parse cell number: %w", err)
	}

	game.stepCellN = cellN

	game.log(game.String())

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

func (game *Game) changeTo(color State) func(i int) {
	return func(i int) {
		game.cells[i] = color
	}
}

func (game *Game) count(cellN int, direction direction, color State, change ...func(i int)) int {
	changeFunc := func(_ int) {}
	if len(change) != 0 {
		changeFunc = change[0]
	}
	count := 0

	switch direction {
	case up:
		i := cellN - 8
		for ; i < 0 || game.cells[i] != color; i -= 8 {
			if i < 8 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case down:
		i := cellN + 8
		for ; i > 63 || game.cells[i] != color; i += 8 {
			if i > 63-8 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	case left:
		i := cellN - 1
		for ; i < 0 || game.cells[i] != color; i-- {
			if i < 0 || i%8 == 0 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case right:
		i := cellN + 1
		for ; i > 63 || game.cells[i] != color; i++ {
			if i > 63 || i%8 == 7 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	case leftup:
		i := cellN - 9
		for ; i < 0 || game.cells[i] != color; i -= 9 {
			if i%8 == 0 || i < 8 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case rightup:
		i := cellN - 7
		for ; i < 0 || game.cells[i] != color; i -= 7 {
			if i%8 == 7 || i < 8 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i < 0 {
			return 0
		}

	case rightdown:
		i := cellN + 9
		for ; i > 63 || game.cells[i] != color; i += 9 {
			if i%8 == 7 || i > 63-8 || game.cells[i] == Empty {
				return 0
			}
			changeFunc(i)
			count++
		}
		if i > 63 {
			return 0
		}

	case leftdown:
		i := cellN + 7
		for ; i > 63 || game.cells[i] != color; i += 7 {
			if i%8 == 0 || i > 63-8 || game.cells[i] == Empty {
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
