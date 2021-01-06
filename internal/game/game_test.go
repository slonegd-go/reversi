package game

import (
	"strings"
	"testing"

	"github.com/slonegd-go/reversi/internal/player"
	"github.com/slonegd-go/reversi/internal/player/cli"
	"github.com/stretchr/testify/assert"
)

func TestGame_Step(t *testing.T) {
	tests := map[string]struct {
		game     *Game
		color    player.Color
		position string
		wantErr  string
		wantGame string
	}{
		"E3 green down must ok": {game: g(""), color: Green, position: "E3", wantGame: `
\ A B C D E F G H
1                
2                
3         G      
4       G G      
5       R G      
6                
7                
8                
`},

		"D3 green down must not ok": {game: g(""), color: Green, position: "D3",
			wantErr: "unavailable step", wantGame: `
\ A B C D E F G H
1                
2                
3       x        
4       G R      
5       R G      
6                
7                
8                
`},
		"several directions": {game: g("D2:Green,E2:Green,F2:Red,D3:Green"), color: Red, position: "C2",
			wantGame: `
\ A B C D E F G H
1                
2     R R R R    
3       R        
4       G R      
5       R G      
6                
7                
8                
`},
	}

	for name, tt := range tests {
		err := tt.game.Step(tt.color, tt.position)
		assert.Equal(t, tt.wantGame, tt.game.String(), name)
		if tt.wantErr != "" {
			ok := assert.Error(t, err, name)
			if !ok {
				continue
			}
			assert.Equal(t, tt.wantErr, err.Error())
			continue
		}
		assert.NoError(t, err, name)
	}
}

func TestGame_count(t *testing.T) {
	tests := map[string]struct {
		name      string
		game      *Game
		cellN     int
		direction direction
		color     player.Color
		want      int
		wantGame  string
	}{
		"↖ border": {game: g(""), cellN: n("A1"), color: Red, direction: up, want: 0},
		"← border": {game: g(""), cellN: n("A1"), color: Red, direction: left, want: 0},
		// up cases
		"↑ bad with empty":  {game: g("B1:Empty,B2:Green,B3:Green"), cellN: n("B4"), color: Red, direction: up, want: 0},
		"↑ bad with border": {game: g("B1:Green,B2:Green,B3:Green"), cellN: n("B4"), color: Red, direction: up, want: 0},
		"↑ good":            {game: g("B1:Red,B2:Green,B3:Green"), cellN: n("B4"), color: Red, direction: up, want: 2},
		// down cases
		"↓ bad with empty":  {game: g("B8:Empty,B7:Green,B6:Green"), cellN: n("B5"), color: Red, direction: down, want: 0},
		"↓ bad with border": {game: g("B8:Green,B7:Green,B6:Green"), cellN: n("B5"), color: Red, direction: down, want: 0},
		"↓ good":            {game: g("B8:Red,B7:Green,B6:Green"), cellN: n("B5"), color: Red, direction: down, want: 2},
		// left cases
		"← bad with empty":  {game: g("A2:Empty,B2:Red,C2:Red"), cellN: n("D2"), color: Green, direction: left, want: 0},
		"← bad with border": {game: g("A2:Red,B2:Red,C2:Red"), cellN: n("D2"), color: Green, direction: left, want: 0},
		"← good":            {game: g("A2:Green,B2:Red,C2:Red"), cellN: n("D2"), color: Green, direction: left, want: 2},
		// right cases
		"→ bad with empty":  {game: g("H6:Empty,G6:Empty,F6:Red"), cellN: n("E6"), color: Green, direction: right, want: 0},
		"→ bad with border": {game: g("H6:Red,G6:Red,F6:Red"), cellN: n("E6"), color: Green, direction: right, want: 0},
		"→ good":            {game: g("G6:Green,F6:Red"), cellN: n("E6"), color: Green, direction: right, want: 1},
		// left up cases
		"↖ good":            {game: g("B1:Green,C2:Red"), cellN: n("D3"), color: Green, direction: leftup, want: 1},
		"↖ bad up border":   {game: g("B1:Red,C2:Red"), cellN: n("D3"), color: Green, direction: leftup, want: 0},
		"↖ bad left border": {game: g("A2:Red,B3:Red"), cellN: n("C4"), color: Green, direction: leftup, want: 0},
		"↖ bad empty":       {game: g("B2:Empty,C3:Red"), cellN: n("D4"), color: Green, direction: leftup, want: 0},
		// right up cases
		"↗ good":             {game: g("G1:Green,F2:Red"), cellN: n("E3"), color: Green, direction: rightup, want: 1},
		"↗ bad up border":    {game: g("G1:Red,F2:Red"), cellN: n("E3"), color: Green, direction: rightup, want: 0},
		"↗ bad right border": {game: g("H2:Red,G3:Red"), cellN: n("F4"), color: Green, direction: rightup, want: 0},
		"↗ bad empty":        {game: g("H2:Empty,G3:Red"), cellN: n("F4"), color: Green, direction: rightup, want: 0},
		// right down cases
		"↘ good":             {game: g("H7:Green,G6:Red"), cellN: n("F5"), color: Green, direction: rightdown, want: 1},
		"↘ bad right border": {game: g("H7:Red,G6:Red"), cellN: n("F5"), color: Green, direction: rightdown, want: 0},
		"↘ bad down border":  {game: g("G8:Red,F7:Red"), cellN: n("E5"), color: Green, direction: rightdown, want: 0},
		// left down cases
		"↙ good":            {game: g("A7:Green,B6:Red"), cellN: n("C5"), color: Green, direction: leftdown, want: 1},
		"↙ bad left border": {game: g("A7:Red,B6:Red"), cellN: n("C5"), color: Green, direction: leftdown, want: 0},
		"↙ bad down border": {game: g("B8:Red,C7:Red"), cellN: n("D6"), color: Green, direction: leftdown, want: 0},
		"↙ bad empty":       {game: g("B8:Empty,C7:Red"), cellN: n("D6"), color: Green, direction: leftdown, want: 0},
	}
	for name, tt := range tests {
		got := tt.game.count(tt.cellN, tt.direction, tt.color)
		assert.Equal(t, tt.want, got, name)
	}
}

//
//
// helpers and mocks
//
//

func cellN(s string) int {
	result, _ := parseCellN(s)
	return result
}

func n(s string) int {
	result, _ := parseCellN(s)
	return result
}

func c(s string) player.Color {
	switch s {
	case "Red":
		return player.Red
	case "Green":
		return player.Green
	}
	return player.Empty
}

func g(description string) *Game {
	result := New(&cli.Player{}, &cli.Player{})
	if description == "" {
		return result
	}
	cells := strings.Split(description, ",")
	for _, cell := range cells {
		nColor := strings.Split(cell, ":")
		result.cells[cellN(nColor[0])] = c(nColor[1])
	}

	return result
}

var (
	Green = player.Green
	Red   = player.Red
)
