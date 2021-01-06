package neural

import (
	"github.com/patrikeh/go-deep"
	"github.com/slonegd-go/reversi/internal/player"
)

type Player struct {
	color   player.Color
	neural  *deep.Neural  // обучение во время игры
	weights [][][]float64 // обучение сохраняется при выйгрыше
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
	}
}

func (player *Player) Step(colors []player.Color) string {
	// преоразовать colors во входные для сети
	// выполнить предикт
	// рассортировать полученные результаты
	// сделать шаг с максимальным по модулю
	// если шаг неудачный, то понизить, удачный повысить
	// если выйграл - сохранить, иначе вернуть до игры
	return ""
}

func (p *Player) Notify(result player.Result) {
	if result == player.Win {
		p.weights = p.neural.Weights()
		// TODO записать веса в файл
		return
	}
	p.neural.ApplyWeights(p.weights) // в случае проигыша не запоминаем обученое
}
func (player *Player) SetColor(v player.Color) { player.color = v }
func (player *Player) Color() player.Color     { return player.color }
