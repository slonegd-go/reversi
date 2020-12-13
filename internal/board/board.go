package board

import (
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
