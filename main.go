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
		log.Println(err)
	}
	fmt.Println(gameBoard)

	err = gameBoard.Step(board.Red, "G7")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(gameBoard)

	fmt.Println(gameBoard.Show("E3"))
	err = gameBoard.Step(board.Green, "E3")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(gameBoard)
}
