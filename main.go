package main

import (
	"log"
	"time"

	"github.com/slonegd-go/reversi/internal/game"
	"github.com/slonegd-go/reversi/internal/player/cli"
	"github.com/slonegd-go/reversi/internal/player/neural"
)

func main() {
	p1 := &cli.Player{}
	p2 := neural.New("")
	currentGame := game.New(p1, p2, game.WithLogger(log.Printf))
	currentGame.Start()

	time.Sleep(1 * time.Second)

	// err := currentGame.Step(game.Green, "B3")
	// if err != nil {
	// 	log.Println(err)
	// }

	// err = currentGame.Step(game.Red, "G7")
	// if err != nil {
	// 	log.Println(err)
	// }

	// err = currentGame.Step(game.Green, "E3")
	// if err != nil {
	// 	log.Println(err)
	// }
}
