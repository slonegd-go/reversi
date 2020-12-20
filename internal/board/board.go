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
	cells     [64]State
	stepCellN int
}

func New() *Board {
	cells := [64]State{}
	cells[27] = Green
	cells[28] = Red
	cells[35] = Red
	cells[36] = Green
	return &Board{
		cells:     cells,
		stepCellN: -1,
	}
}

func (board Board) String() string {
	var builder strings.Builder
	builder.Grow(101)
	builder.WriteString("\n\\ A B C D E F G H\n")
	for i := 0; i < 8; i++ {
		builder.WriteString(strconv.Itoa(i + 1))
		for j := 0; j < 8; j++ {
			number := i*8 + j
			if board.stepCellN == number {
				builder.WriteString(" x")
				continue
			}
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

func (board *Board) Step(color State, position string) error {
	if color != Green && color != Red {
		return errors.New("only green and red state available")
	}

	cellN, err := parseCellN(position)
	if err != nil {
		return fmt.Errorf("parse cell number: %w", err)
	}

	board.stepCellN = cellN

	directions, err := board.check(int(cellN), color)
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}
	board.cells[cellN] = color
	board.fill(cellN, directions, color)
	return nil
}

func (board Board) check(cellN int, color State) ([]direction, error) {
	if board.cells[cellN] != Empty {
		return nil, errors.New("cell not empty")
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
		return nil, errors.New("no other color beside")
	}

	return directions, nil
}

func (board *Board) fill(cellN int, directions []direction, color State) int {
	count := 1
	board.cells[cellN] = color
	for _, direction := range directions {
		switch direction {
		case up:
			for i := cellN - 8; board.cells[i].not(color); i -= 8 {
				board.cells[i] = color
				count++
			}

		case down:
			for i := cellN + 8; board.cells[i].not(color); i += 8 {
				board.cells[i] = color
				count++
			}

		case left:
			for i := cellN - 1; board.cells[i].not(color); i-- {
				board.cells[i] = color
				count++
			}

		case right:
			for i := cellN + 1; board.cells[i].not(color); i++ {
				board.cells[i] = color
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
