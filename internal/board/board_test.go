package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_Step(t *testing.T) {
	tests := []struct {
		name     string
		board    *Board
		color    State
		position string
		wantErr  string
	}{
		{
			name:     "E3 green down must ok",
			board:    New(),
			color:    Green,
			position: "E3",
		},
		{
			name:     "D3 green down must not ok",
			board:    New(),
			color:    Green,
			position: "D3",
			wantErr:  "check: no other color beside",
		},
	}
	for _, tt := range tests {
		err := tt.board.Step(tt.color, tt.position)
		if tt.wantErr != "" {
			assert.Error(t, err, tt.name)
			assert.Equal(t, tt.wantErr, err.Error())
			continue
		}
		assert.NoError(t, err, tt.name)
	}
}

func TestBoard_fill(t *testing.T) {
	tests := map[string]struct {
		name       string
		board      *Board
		cellN      int
		directions []direction
		color      State
		want       int
		wantBoard  string
	}{
		"up": {
			board:      New(),
			directions: []direction{up},
			cellN:      cellN("D6"),
			color:      Green,
			want:       2,
			wantBoard: `
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
			board:      New(),
			directions: []direction{down},
			cellN:      cellN("D3"),
			color:      Red,
			want:       2,
			wantBoard: `
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
			board:      New(),
			directions: []direction{right},
			cellN:      cellN("C4"),
			color:      Red,
			want:       2,
			wantBoard: `
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
			board:      New(),
			directions: []direction{left},
			cellN:      cellN("F4"),
			color:      Green,
			want:       2,
			wantBoard: `
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
			board: func() *Board {
				board := New()
				board.cells[cellN("F4")] = Green
				board.cells[cellN("F3")] = Red
				return board
			}(),
			directions: []direction{left, up},
			cellN:      cellN("F5"),
			color:      Red,
			want:       3,
			wantBoard: `
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
		got := tt.board.fill(tt.cellN, tt.directions, tt.color)
		assert.Equal(t, tt.want, got, name)
		assert.Equal(t, tt.wantBoard, tt.board.String(), name)
	}
}

func cellN(s string) int {
	result, _ := parseCellN(s)
	return result
}
