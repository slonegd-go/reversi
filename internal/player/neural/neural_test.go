package neural

import (
	"testing"

	"github.com/slonegd-go/reversi/internal/player"
	"github.com/stretchr/testify/assert"
)

func TestPlayer_updateInputs(t *testing.T) {
	tests := []struct {
		name   string
		colors []player.Color
		color  player.Color
		want   []float64
	}{
		{
			name:   "our",
			colors: append([]player.Color{g}, generateColors(e, 63)...),
			color:  g,
			want:   append([]float64{0xFFFFFFFD, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF}, generateFloats(0, 116)...),
		},
		{
			name:   "not our",
			colors: append([]player.Color{r}, generateColors(e, 63)...),
			color:  g,
			want:   append([]float64{0xFFFFFFFE, 0xFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF}, generateFloats(0, 116)...),
		},
	}
	for _, tt := range tests {
		p := New("")
		p.SetColor(tt.color)
		p.updateInputs(tt.colors)
		assert.Equal(t, tt.want, p.inputs, tt.name)
	}
}

func Test_cell(t *testing.T) {
	tests := map[string]struct {
		i    int
		want string
	}{
		"A1": {i: 0, want: "A1"},
		"H1": {i: 7, want: "H1"},
		"C4": {i: 24 + 3 - 1, want: "C4"},
		"F4": {i: 24 + 4 - 1, want: "F4"},
		"C5": {i: 32 + 3 - 1 - 2, want: "C5"},
		"F5": {i: 32 + 4 - 1 - 2, want: "F5"},
	}
	for name, tt := range tests {
		assert.Equal(t, tt.want, cell(tt.i), name)
	}
}

//
//
// helpers and mocks
//
//

var (
	g = player.Green
	r = player.Red
	e = player.Empty
)

func generateColors(color player.Color, count int) []player.Color {
	result := []player.Color{}
	for i := 0; i < count; i++ {
		result = append(result, color)
	}
	return result
}

func generateFloats(f float64, count int) []float64 {
	result := []float64{}
	for i := 0; i < count; i++ {
		result = append(result, f)
	}
	return result
}
