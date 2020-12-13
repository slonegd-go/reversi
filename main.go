package main

import (
	"fmt"
	"log"

	"github.com/slonegd-go/reversi/internal/board"
)

func main() {
	gameBoard := board.New()
	fmt.Println(gameBoard)

	err := gameBoard.Step(board.Green, "B3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gameBoard)

	gameBoard.Step(board.Red, "G7")
	fmt.Println(gameBoard)
}
