package neural

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
	"github.com/slonegd-go/reversi/internal/player"
)

type Player struct {
	color    player.Color
	neural   *deep.Neural // обучение во время игры
	persist  persist
	inputs   []float64
	index    int
	trainer  training.Trainer
	filename string
}

type persist struct {
	Weights             [][][]float64
	WinCount, LoseCount int
}

func (p *Player) Stats() {
	log.Printf("%s win\t%d:%d\tlose, ratio %f", p.filename, p.persist.WinCount, p.persist.LoseCount, float32(p.persist.WinCount)/float32(p.persist.LoseCount))
}

func New(filename string) *Player {
	neural := deep.NewNeural(&deep.Config{
		Inputs:     240,
		Layout:     []int{360, 480, 600, 480, 360, 240, 120, 60},
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeRegression,
		Weight:     deep.NewNormal(1, 0),
		Bias:       true,
	})
	persist := persist{
		Weights: neural.Weights(),
	}

	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		gob.NewDecoder(file).Decode(&persist)
		neural.ApplyWeights(persist.Weights)
	}

	return &Player{
		neural:   neural,
		persist:  persist,
		inputs:   make([]float64, 240),
		trainer:  training.NewTrainer(training.NewSGD(0.005, 0.5, 1e-6, true), 1),
		filename: filename,
	}
}

func (p *Player) Step(colors []player.Color, enabledCells []bool, step func(string) error) {

	p.updateInputs(colors)

	for {
		// time.Sleep(200 * time.Millisecond)
		// log.Printf("inputs: %+v", p.inputs)
		outputs := p.neural.Predict(p.inputs)

		predict := []output{}
		for i, f64 := range outputs {
			n := cellN(i)
			if !enabledCells[n] {
				outputs[i] = 0
				f64 = 0
			}

			predict = append(predict, output{
				i:      i,
				cell:   cell(n),
				weight: abs(f64),
			})

		}
		sort.Slice(predict, func(i, j int) bool {
			return predict[i].weight > predict[j].weight // по убыванию
		})
		// log.Printf("predict: %+v", predict)
		log.Printf("predict:\n\t%+v,\n\t%+v,\n\t%+v,\n\t%+v,\n\t%+v",
			predict[0], predict[1], predict[2], predict[3], predict[4])

		k := 10.
		err := step(predict[0].cell)
		if err != nil {
			k = 0
		}

		outputs[predict[0].i] *= k

		// TODO если проиграл, то понизить кофжжфициенты, которые к этому привели
		// тренировать только один раз - по результату игры
		p.trainer.Train(p.neural, training.Examples{
			{
				Input:    p.inputs,
				Response: outputs,
			},
		}, nil, 1)

		if err == nil {
			return
		}
	}
}

func (p *Player) Notify(result player.Result) {
	p.inputs = make([]float64, 240)
	p.index = 0

	if result == player.Win {
		p.persist.Weights = p.neural.Weights()
		p.persist.WinCount++
	} else {
		p.persist.LoseCount++
		p.neural.ApplyWeights(p.persist.Weights) // в случае проигыша не запоминаем обученое
	}

	file, err := os.Create(p.filename)
	if err == nil {
		gob.NewEncoder(file).Encode(p.persist)
	}
}

func (player *Player) SetColor(v player.Color) { player.color = v }
func (player *Player) Color() player.Color     { return player.color }

func (p *Player) updateInputs(colors []player.Color) {
	uints := [8]uint16{}
	for i, color := range colors {
		offset := (i % 8) * 2
		index := i / 8
		if color == player.Empty {
			continue
		}
		if color == p.color {
			uints[index] = uints[index] | 0b01<<offset
			continue
		}
		uints[index] = uints[index] | 0b10<<offset
	}
	for _, ui := range uints {
		f := 0.
		if ui != 0 {
			f = 1. / float64(ui)
		}
		p.inputs[p.index] = f
		p.index++
	}
}

type output struct {
	i      int
	cell   string
	weight float64
}

func abs(v float64) float64 {
	if v >= 0 {
		return v
	}
	return -v
}

func cellN(i int) int {
	if i > 24+3-1 { // D4, E4 пропускаем
		i += 2
	}
	if i > 32+3-1 { // D5, E5 пропускаем
		i += 2
	}
	return i
}

func cell(i int) string {
	letter := rune('A' + i%8)
	index := rune('1' + i/8)
	return fmt.Sprintf("%c%c", letter, index)
}
