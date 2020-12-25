package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	file, err := os.Open("./F_C.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	priceDataChan := make(chan StockData)
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
		close(priceDataChan)
	}()

	upChan := make(chan StockData)
	go func() {
		for price := range priceDataChan {
			if price.close > price.open {
				upChan <- price
			}
		}
		close(upChan)
	}()

	for price := range upChan {
		fmt.Printf("%s, open %f, close %f, diff %f\n",
			price.date.String(),
			price.open,
			price.close,
			price.close-price.open)
	}
}
