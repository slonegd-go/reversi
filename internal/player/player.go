package player

import "fmt"

type Color int
type Result int

const (
	Empty Color = iota
	Green
	Red
)

func (c Color) String() string {
	switch c {
	case Empty:
		return "empty"
	case Green:
		return "green"
	case Red:
		return "red"
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
