package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/slonegd-go/reversi/internal/evolution"
	"github.com/slonegd-go/reversi/internal/game"
	"github.com/slonegd-go/reversi/internal/player/cli"
	"github.com/slonegd-go/reversi/internal/player/neural"
)

func main() {

	stats := flag.Int("stats", 0, "return stats of epoch")
	player := flag.String("player", "", "play with neural")
	flag.Parse()

	if *stats != 0 {
		epoch := *stats
		path := filepath.Join(".", "players", fmt.Sprintf("epoch%d", epoch))
		players := []*neural.Player{
			neural.New(path, fmt.Sprintf("%d_1", epoch)),
			neural.New(path, fmt.Sprintf("%d_2", epoch)),
			neural.New(path, fmt.Sprintf("%d_3", epoch)),
			neural.New(path, fmt.Sprintf("%d_4", epoch)),
			neural.New(path, fmt.Sprintf("%d_5", epoch)),
			neural.New(path, fmt.Sprintf("%d_6", epoch)),
			neural.New(path, fmt.Sprintf("%d_7", epoch)),
			neural.New(path, fmt.Sprintf("%d_8", epoch)),
			neural.New(path, fmt.Sprintf("%d_9", epoch)),
		}
		for _, p := range players {
			if p.WinCount() != 0 {
				p.Stats()
			}
		}

		gameCount := 0
		for _, player := range players {
			gameCount += player.WinCount()
		}
		log.Printf("games count %d", gameCount)
		return
	}

	if *player != "" {
		tmp := strings.Split(*player, "_")
		epoch := tmp[0]
		path := filepath.Join(".", "players", fmt.Sprintf("epoch%s", epoch))
		n := neural.New(path, *player)
		p := &cli.Player{}
		currentGame := game.New(n, p, game.WithLogger(log.Printf))
		currentGame.Start()
		time.Sleep(1 * time.Second)
		return
	}

	rand.Seed(time.Now().UnixNano())
	evolution.Start()

}
