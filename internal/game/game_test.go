package game

import (
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

func cellN(s string) int {
	result, _ := parseCellN(s)
	return result
}
