package main

import (
	"fmt"

	"github.com/slonegd-go/reversi/internal/board"
)

func main() {
	board := board.New()

	fmt.Println(board)
}
