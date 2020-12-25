package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	done := make(chan int)
	fmt.Println(files)
	// Stage1 - Process files in parallel
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		go func() {
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
			done <- 1
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
		var max StockData
		for new := range upChan {
			fmt.Println(max.close - max.open)
			fmt.Println(new.close - new.open)

			if max.close-max.open < new.close-new.open {
				max = new
			}

			fmt.Println(max)
		}
		maxChan <- max
	}()

	// Wait until stage 1 is done
	for range files {
		<-done
	}
	close(priceDataChan)
	close(done)

	max := <-maxChan
	fmt.Printf("%s, open %f, close %f, diff %f\n",
		max.date.String(),
		max.open,
		max.close,
		max.close-max.open)

}
