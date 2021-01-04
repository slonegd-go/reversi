package game

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_Step(t *testing.T) {
	tests := []struct {
		name     string
		game     *Game
		color    State
		position string
		wantErr  string
	}{
		{
			name:     "E3 green down must ok",
			game:     New(),
			color:    Green,
			position: "E3",
		},
		{
			name:     "D3 green down must not ok",
			game:     New(),
			color:    Green,
			position: "D3",
			wantErr:  "check: no other color beside",
		},
	}
	for _, tt := range tests {
		err := tt.game.Step(tt.color, tt.position)
		if tt.wantErr != "" {
			assert.Error(t, err, tt.name)
			assert.Equal(t, tt.wantErr, err.Error())
			continue
		}
		assert.NoError(t, err, tt.name)
	}
}

func TestGame_fill(t *testing.T) {
	tests := map[string]struct {
		name       string
		game       *Game
		cellN      int
		directions []direction
		color      State
		want       int
		wantGame   string
	}{
		"up": {
			game:       New(),
			directions: []direction{up},
			cellN:      cellN("D6"),
			color:      Green,
			want:       2,
			wantGame: `
\ A B C D E F G H
1                
2                
3                
4       G R      
5       G G      
6       G        
7                
8                
`,
		},

		"down": {
			game:       New(),
			directions: []direction{down},
			cellN:      cellN("D3"),
			color:      Red,
			want:       2,
			wantGame: `
\ A B C D E F G H
1                
2                
3       R        
4       R R      
5       R G      
6                
7                
8                
`,
		},

		"right": {
			game:       New(),
			directions: []direction{right},
			cellN:      cellN("C4"),
			color:      Red,
			want:       2,
			wantGame: `
\ A B C D E F G H
1                
2                
3                
4     R R R      
5       R G      
6                
7                
8                
`,
		},

		"left": {
			game:       New(),
			directions: []direction{left},
			cellN:      cellN("F4"),
			color:      Green,
			want:       2,
			wantGame: `
\ A B C D E F G H
1                
2                
3                
4       G G G    
5       R G      
6                
7                
8                
`,
		},

		"left and up": {
			game: func() *Game {
				game := New()
				game.cells[cellN("F4")] = Green
				game.cells[cellN("F3")] = Red
				return game
			}(),
			directions: []direction{left, up},
			cellN:      cellN("F5"),
			color:      Red,
			want:       3,
			wantGame: `
\ A B C D E F G H
1                
2                
3           R    
4       G R R    
5       R R R    
6                
7                
8                
`,
		},
	}
	for name, tt := range tests {
		got := tt.game.fill(tt.cellN, tt.directions, tt.color)
		assert.Equal(t, tt.want, got, name)
		assert.Equal(t, tt.wantGame, tt.game.String(), name)
	}
}

func TestGame_count(t *testing.T) {
	tests := map[string]struct {
		name      string
		game      *Game
		cellN     int
		direction direction
		color     State
		want      int
		wantGame  string
	}{
		"↖ border": {game: New(), cellN: n("A1"), color: Red, direction: up, want: 0},
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

func color(s string) State {
	switch s {
	case "Red":
		return Red
	case "Green":
		return Green
	}
	return Empty
}

func g(description string) *Game {
	result := New()
	cells := strings.Split(description, ",")
	for _, cell := range cells {
		nColor := strings.Split(cell, ":")
		result.cells[cellN(nColor[0])] = color(nColor[1])
	}

	return result
}
