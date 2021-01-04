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

func (player *Player) Step([]player.Color) string {
	reader := bufio.NewReader(os.Stdin)
	result, _ := reader.ReadString('\n')
	return strings.TrimSpace(result)
}
func (player *Player) Notify(player.Result)    {}
func (player *Player) SetColor(v player.Color) { player.color = v }
func (player *Player) Color() player.Color     { return player.color }
