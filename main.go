package main

import (
	"log"

	"github.com/slonegd-go/reversi/internal/game"
)

func main() {
	currentGame := game.New(game.WithLogger(log.Printf))

	err := currentGame.Step(game.Green, "B3")
	if err != nil {
		log.Println(err)
	}

	err = currentGame.Step(game.Red, "G7")
	if err != nil {
		log.Println(err)
	}

	err = currentGame.Step(game.Green, "E3")
	if err != nil {
		log.Println(err)
	}
}
