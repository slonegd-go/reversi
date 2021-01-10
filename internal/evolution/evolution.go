package evolution

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/slonegd-go/reversi/internal/game"
	"github.com/slonegd-go/reversi/internal/player/neural"
)

func Start() {

	for epoch := 1; ; epoch++ {
		log.Printf("start epoch #%d", epoch)
		// определить эпоху
		path := filepath.Join(".", "players", fmt.Sprintf("epoch%d", epoch+1))
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		// загрузить
		path = filepath.Join(".", "players", fmt.Sprintf("epoch%d", epoch))
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

		// определить сколько игр прошло
		gameCount := 0
		for _, player := range players {
			gameCount += player.WinCount()
		}

		// продолжить обучение
		for ; gameCount < 10000; gameCount += 4 {
			log.Printf("start %d game of %d epoch", gameCount, epoch)
			plN := make([]int, 0, 8)
			plN = append(plN, rand.Intn(len(players)))
			for i := 1; i < 8; i++ {
				for {
					n := rand.Intn(len(players))
					if exist(plN, n) {
						continue
					}
					plN = append(plN, n)
					break
				}
			}

			var wg sync.WaitGroup
			wg.Add(4)
			for i := 0; i < 8; i += 2 {
				go func(i int) {
					currentGame := game.New(players[plN[i]], players[plN[i+1]], game.WithLogger(log.Printf))
					currentGame.Start()
					wg.Done()
				}(i)
			}
			wg.Wait()
		}

		// по окончанию определить лучших
		sort.Slice(players, func(i, j int) bool {
			return players[i].WinRatio() > players[j].WinRatio()
		})

		// сгенерировать новых
		newEpoch := epoch + 1
		regenerate := true
		path = filepath.Join(".", "players", fmt.Sprintf("epoch%d", newEpoch))
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Printf(err.Error())
			return
		}

		players[0].CopyToFilename(path, fmt.Sprintf("%d_1", newEpoch))
		players[0].CopyToFilename(path, fmt.Sprintf("%d_2", newEpoch), regenerate)
		players[0].CopyToFilename(path, fmt.Sprintf("%d_3", newEpoch), regenerate)
		players[1].CopyToFilename(path, fmt.Sprintf("%d_4", newEpoch))
		players[1].CopyToFilename(path, fmt.Sprintf("%d_5", newEpoch), regenerate)
		players[1].CopyToFilename(path, fmt.Sprintf("%d_6", newEpoch), regenerate)
		players[2].CopyToFilename(path, fmt.Sprintf("%d_7", newEpoch))
		players[2].CopyToFilename(path, fmt.Sprintf("%d_8", newEpoch), regenerate)
		players[2].CopyToFilename(path, fmt.Sprintf("%d_9", newEpoch), regenerate)
	}
}

func exist(list []int, v int) bool {
	for _, v1 := range list {
		if v1 == v {
			return true
		}
	}
	return false
}
