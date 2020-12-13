package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.Create("./data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"minN", "maxN", "d1", "d2", "d3", "d4", "d5"})
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	for count := 0; count < 100; count++ {
		data := make([]int, 0, 5)

		minN, maxN := 0, 0
		for i := 0; i < 5; i++ {
			data = append(data, rand.Intn(10))
			if data[i] < data[minN] {
				minN = i
			}
			if data[i] > data[maxN] {
				maxN = i
			}
		}
		log.Println("data", data)
		log.Println("minN", minN)
		log.Println("maxN", maxN)

		line := make([]string, 0, 7)
		line = append(line, strconv.Itoa(minN))
		line = append(line, strconv.Itoa(maxN))
		for _, d := range data {
			line = append(line, strconv.Itoa(d))
		}

		log.Println(line)
		writer.Write(line)
	}

}
