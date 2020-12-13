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
