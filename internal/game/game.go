package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func New(opts ...Option) *Game {
	cells := [64]State{}
	cells[27] = Green
	cells[28] = Red
	cells[35] = Red
	cells[36] = Green

	options := &Options{
		log: func(string, ...interface{}) {},
	}

	for _, opt := range opts {
		opt(options)
	}

	game := &Game{
		cells:     cells,
		stepCellN: -1,
		log:       options.log,
	}

	game.log(game.String())

	return game
}

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
				builder.WriteString(" G")
			case Red:
				builder.WriteString(" R")
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

	directions, err := game.check(int(cellN), color)
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}
	game.cells[cellN] = color
	game.fill(cellN, directions, color)

	game.log(game.String())
	return nil
}

func (game Game) check(cellN int, color State) ([]direction, error) {
	if game.cells[cellN] != Empty {
		return nil, errors.New("cell not empty")
	}

	column := cellN % 8
	line := cellN / 8

	directions := make([]direction, 0)

	if column != 0 && game.cells[cellN-1].not(color) {
		directions = append(directions, left)
	}

	if column != 7 && game.cells[cellN-1].not(color) {
		directions = append(directions, right)
	}

	if line != 0 && game.cells[cellN-8].not(color) {
		directions = append(directions, up)
	}

	if line != 7 && game.cells[cellN+8].not(color) {
		directions = append(directions, down)
	}

	// TODO diagonals
	// TODO fill with numbers

	if len(directions) == 0 {
		return nil, errors.New("no other color beside")
	}

	return directions, nil
}

func (game *Game) fill(cellN int, directions []direction, color State) int {
	game.stepCellN = -1
	count := 1
	game.cells[cellN] = color
	for _, direction := range directions {
		switch direction {
		case up:
			for i := cellN - 8; game.cells[i].not(color); i -= 8 {
				game.cells[i] = color
				count++
			}

		case down:
			for i := cellN + 8; game.cells[i].not(color); i += 8 {
				game.cells[i] = color
				count++
			}

		case left:
			for i := cellN - 1; game.cells[i].not(color); i-- {
				game.cells[i] = color
				count++
			}

		case right:
			for i := cellN + 1; game.cells[i].not(color); i++ {
				game.cells[i] = color
				count++
			}
		} // switch direction
	} // range directions

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
