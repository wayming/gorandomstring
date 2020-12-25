package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StockData struct {
	date  time.Time
	open  float64
	high  float64
	low   float64
	close float64
}

func main() {

	// files, err := filepath.Glob("../../../tickerData/F_C.txt")
	files, err := filepath.Glob("../../../tickerData/*txt")
	if err != nil {
		panic(err)
	}
	priceDataChan := make(chan StockData, 10240)
	var wg sync.WaitGroup
	fmt.Println(files)
	// Stage1 - Process files in parallel
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader := bufio.NewReader(file)
			for {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				}
				tokens := strings.Split(string(line), ",")
				priceDate, _ := time.Parse("20060102", tokens[0])
				priceOpen, _ := strconv.ParseFloat(tokens[1], 64)
				priceHigh, _ := strconv.ParseFloat(tokens[2], 64)
				priceLow, _ := strconv.ParseFloat(tokens[3], 64)
				priceClose, _ := strconv.ParseFloat(tokens[4], 64)
				data := StockData{priceDate, priceOpen, priceHigh, priceLow, priceClose}
				priceDataChan <- data
			}
		}()
	}

	// Stage 2 - filter the stock that goes up
	upChan := make(chan StockData, 10240)
	go func() {
		for price := range priceDataChan {
			if price.close > price.open {
				upChan <- price
			}
		}
		close(upChan)
	}()

	// Stage 3 - display greatest diff
	maxChan := make(chan StockData)
	go func() {
		max := <-upChan
		for new := range upChan {
			// fmt.Println(max.close / max.open)

			if max.open == 0 {
				max = new
				continue
			}

			if new.open == 0 {
				continue
			}

			if (max.close / max.open) < (new.close / new.open) {
				max = new
				continue
			}

			// fmt.Println(max.close / max.open)
		}
		maxChan <- max
	}()

	// Wait until stage 1 is done
	go func() {
		wg.Wait()
		close(priceDataChan)
	}()

	max := <-maxChan
	fmt.Printf("%s, open %f, close %f, diff %f\n",
		max.date.String(),
		max.open,
		max.close,
		max.close/max.open)

}
