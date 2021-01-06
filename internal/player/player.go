package player

import (
	"fmt"

	"github.com/fatih/color"
)

type Color int
type Result int

const (
	Empty Color = iota
	Green
	Red
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func (c Color) String() string {
	switch c {
	case Empty:
		return "empty"
	case Green:
		return green("green")
	case Red:
		return red("red")
	default:
		return fmt.Sprintf("undefined(%d)", c)
	}
}

const (
	Lose Result = iota
	Win
)

type Player interface {
	Step([]Color) string
	Notify(Result)
	SetColor(Color)
	Color() Color
}
