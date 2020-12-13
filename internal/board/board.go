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

func (board *Board) Step(color State, position string) error {
	if color != Green && color != Red {
		return errors.New("only green and red state available")
	}
	if len(position) != 2 {
		return fmt.Errorf("position must be from A1 to H8, got: %s", position)
	}
	column := position[0] - byte('A')
	line := position[1] - byte('1')
	if column > 7 || line > 7 {
		return fmt.Errorf("position must be from A1 to H8, got: %s", position)
	}

	// TODO check
	board.cells[column+line*8] = color
	return nil
}
