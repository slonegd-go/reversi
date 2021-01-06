package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/slonegd-go/reversi/internal/game"
	"github.com/slonegd-go/reversi/internal/player/neural"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// p1 := &cli.Player{}
	p1 := neural.New("")
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
