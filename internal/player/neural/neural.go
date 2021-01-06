package neural

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/patrikeh/go-deep"
	"github.com/slonegd-go/reversi/internal/player"
)

type Player struct {
	color   player.Color
	neural  *deep.Neural  // обучение во время игры
	weights [][][]float64 // обучение сохраняется при выйгрыше
	inputs  []float64     // 30 ходов, по 64*2 бита на клетку=128 бит на ход, в один float64 влезает 32 бита, итого 4 float на ход
	index   int
}

// TODO прочитать веса из файла
func New(filename string) *Player {
	neural := deep.NewNeural(&deep.Config{
		Inputs:     120,
		Layout:     []int{240, 360, 240, 120, 60},
		Activation: deep.ActivationTanh,
		Mode:       deep.ModeRegression,
		Weight:     deep.NewNormal(1, 0),
		Bias:       true,
	})
	return &Player{
		neural:  neural,
		weights: neural.Weights(),
		inputs:  make([]float64, 120),
	}
}

func (p *Player) Step(colors []player.Color, step func(string) error) {

	p.updateInputs(colors)

	for {
		time.Sleep(500 * time.Millisecond)
		log.Printf("inputs: %+v", p.inputs)
		outputs := p.neural.Predict(p.inputs)

		predict := []output{}
		for i, f64 := range outputs {
			predict = append(predict, output{
				i:      i,
				cell:   cell(i),
				weight: f64,
			})
		}
		sort.Slice(predict, func(i, j int) bool {
			return abs(predict[i].weight) > abs(predict[j].weight) // по убыванию
		})
		log.Printf("outputs: %+v", predict)
		err := step(predict[0].cell)
		if err != nil {
			continue
		}
	}

	// если шаг неудачный, то понизить, удачный повысить
	// если выйграл - сохранить, иначе вернуть до игры

}

func (p *Player) Notify(result player.Result) {
	p.inputs = make([]float64, 120)
	p.index = 0

	if result == player.Win {
		p.weights = p.neural.Weights()
		// TODO записать веса в файл
		return
	}
	p.neural.ApplyWeights(p.weights) // в случае проигыша не запоминаем обученое
}
func (player *Player) SetColor(v player.Color) { player.color = v }
func (player *Player) Color() player.Color     { return player.color }

func (p *Player) updateInputs(colors []player.Color) {
	uints := [4]uint32{}
	for i, color := range colors {
		offset := (i % 16) * 2
		index := i / 16
		if color == player.Empty {
			uints[index] = uints[index] | 0b11<<offset
			continue
		}
		if color == p.color {
			uints[index] = uints[index] | 0b01<<offset
			continue
		}
		uints[index] = uints[index] | 0b10<<offset
	}
	for _, ui := range uints {
		p.inputs[p.index] = 1. / float64(ui)
		p.index++
	}
}

type output struct {
	i      int
	cell   string
	weight float64
}

func abs(v float64) float64 {
	if v > 0 {
		return v
	}
	return -v
}

func cell(i int) string {
	if i > 24+3-1 { // D4, E4 пропускаем
		i += 2
	}
	if i > 32+3-1 { // D5, E5 пропускаем
		i += 2
	}
	letter := rune('A' + i%8)
	index := rune('1' + i/8)
	return fmt.Sprintf("%c%c", letter, index)
}
