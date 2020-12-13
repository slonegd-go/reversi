package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	data, err := load("./data.csv")
	if err != nil {
		log.Fatal(err)
	}

	test, err := load("./test.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("have %d entries\n", len(data))

	neural := deep.NewNeural(&deep.Config{
		Inputs:     len(data[0].Input),
		Layout:     []int{10, 20, 10, 2},
		Activation: deep.ActivationTanh,
		Mode:       deep.ModeRegression,
		Weight:     deep.NewNormal(1, 0),
		Bias:       true,
	})

	got := neural.Predict([]float64{2, 6, 7, 8, 1})
	log.Printf("[2 6 7 8 1] got %v, want [4 3]", got)

	got = neural.Predict([]float64{1, 7, 3, 5, 1})
	log.Printf("[1 5 3 7 1] got %v, want [0 1]", got)

	got = neural.Predict([]float64{5, 2, 8, 3, 4})
	log.Printf("[5 2 8 3 4] got %v, want [1 2]", got)

	const (
		epochs = 50
		count  = 20
	)

	trainer := training.NewTrainer(training.NewSGD(0.005, 0.5, 1e-6, true), epochs)

	for i := 0; i < count; i++ {
		trainer.Train(neural, data, test, epochs)

		got = neural.Predict([]float64{2, 6, 7, 8, 1})
		log.Printf("[2 6 7 8 1] got %v, want [4 3]", got)

		got = neural.Predict([]float64{1, 7, 3, 5, 1})
		log.Printf("[1 5 3 7 1] got %v, want [0 1]", got)

		got = neural.Predict([]float64{5, 2, 8, 3, 4})
		log.Printf("[5 2 8 3 4] got %v, want [1 2]", got)
	}

}

func load(path string) (training.Examples, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bufio.NewReader(f))

	var examples training.Examples
	for i := 0; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if i == 0 { // header
			continue
		}
		examples = append(examples, example(record))
	}

	return examples, nil
}

func example(in []string) training.Example {
	resEncoded := make([]float64, 0, 2)
	for i := 0; i < 2; i++ {
		res, err := strconv.ParseFloat(in[i], 64)
		if err != nil {
			log.Fatal(err)
		}
		resEncoded = append(resEncoded, res)
	}

	var features []float64
	for i := 2; i < len(in); i++ {
		res, err := strconv.ParseFloat(in[i], 64)
		if err != nil {
			log.Fatal(err)
		}
		features = append(features, res)
	}

	return training.Example{
		Response: resEncoded,
		Input:    features,
	}
}
