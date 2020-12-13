package board

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

type Board struct {
	cells [64]State
}

func New() *Board {
	cells := [64]State{}
	cells[27] = Green
	cells[28] = Red
	cells[35] = Red
	cells[36] = Green
	return &Board{
		cells: cells,
	}
}

func (board Board) String() string {
	var builder strings.Builder
	builder.Grow(101)
	builder.WriteString("\n\\ A B C D E F G H\n")
	for i := 0; i < 8; i++ {
		builder.WriteString(strconv.Itoa(i + 1))
		for j := 0; j < 8; j++ {
			switch board.cells[i*8+j] {
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

func (board Board) Show(position string) string {
	cellN, err := parseCellN(position)
	if err != nil {
		cellN = -1
	}

	var builder strings.Builder
	builder.Grow(101)
	builder.WriteString("\n\\ A B C D E F G H\n")
	for i := 0; i < 8; i++ {
		builder.WriteString(strconv.Itoa(i + 1))
		for j := 0; j < 8; j++ {
			number := i*8 + j
			if cellN == number {
				builder.WriteString(" x")
				continue
			}
			switch board.cells[number] {
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

func (board *Board) Step(color State, position string) error {
	if color != Green && color != Red {
		return errors.New("only green and red state available")
	}

	cellN, err := parseCellN(position)
	if err != nil {
		return fmt.Errorf("parse cell number: %w", err)
	}

	err = board.check(int(cellN), color)
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}
	board.cells[cellN] = color
	return nil
}

func (board Board) check(cellN int, color State) error {
	if board.cells[cellN] != Empty {
		return errors.New("cell not empty")
	}

	column := cellN % 8
	line := cellN / 8

	directions := make([]direction, 0)

	if column != 0 && board.cells[cellN-1].not(color) {
		directions = append(directions, left)
	}

	if column != 7 && board.cells[cellN-1].not(color) {
		directions = append(directions, right)
	}

	if line != 0 && board.cells[cellN-8].not(color) {
		directions = append(directions, up)
	}

	if line != 7 && board.cells[cellN+8].not(color) {
		directions = append(directions, down)
	}

	// TODO diagonals
	// TODO fill with numbers

	if len(directions) == 0 {
		return errors.New("no other color beside")
	}

	return nil
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
