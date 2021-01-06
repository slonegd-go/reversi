package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/slonegd-go/reversi/internal/game"
	"github.com/slonegd-go/reversi/internal/player/cli"
	"github.com/slonegd-go/reversi/internal/player/neural"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	stats := flag.Bool("stats", false, "return stats")
	player := flag.Int("player", 0, "play with neural")
	flag.Parse()

	// p1 := &cli.Player{}
	players := []*neural.Player{
		neural.New("./players/01.go"),
		neural.New("./players/02.go"),
		neural.New("./players/03.go"),
		neural.New("./players/04.go"),
		neural.New("./players/05.go"),
		neural.New("./players/06.go"),
		neural.New("./players/07.go"),
		neural.New("./players/08.go"),
		neural.New("./players/09.go"),
		neural.New("./players/10.go"),
	}
	if *stats {
		for _, p := range players {
			p.Stats()
		}
		return
	}
	if *player > 0 {
		p := &cli.Player{}
		currentGame := game.New(players[*player-1], p, game.WithLogger(log.Printf))
		currentGame.Start()
		time.Sleep(100 * time.Millisecond)
		return
	}

	for {
		p1 := rand.Intn(10)
		p2 := rand.Intn(10)
		for p1 == p2 {
			p2 = rand.Intn(10)
		}
		currentGame := game.New(players[p1], players[p2], game.WithLogger(log.Printf))

		currentGame.Start()

		log.Printf("green is %d, red is %d", p1, p2)
		// time.Sleep(10 * time.Second)
	}

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
