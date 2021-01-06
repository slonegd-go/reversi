package cli

import (
	"bufio"
	"os"
	"strings"

	"github.com/slonegd-go/reversi/internal/player"
)

type Player struct {
	color player.Color
}

func (player *Player) Step(colors []player.Color, step func(string) error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		result, _ := reader.ReadString('\n')
		err := step(strings.TrimSpace(result))
		if err == nil {
			return
		}
	}
}
func (player *Player) Notify(player.Result)    {}
func (player *Player) SetColor(v player.Color) { player.color = v }
func (player *Player) Color() player.Color     { return player.color }
