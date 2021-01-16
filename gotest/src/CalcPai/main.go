package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	//totalIterations := 100000000000
	totalIterations := 1000000000
	workerIterations := 100000000

	countChan := make(chan int)
	numThread := 0
	leftIterations := totalIterations
	rand.Seed(time.Now().UnixNano())

	for {
		chunk := int(math.Min(float64(leftIterations), float64(workerIterations)))

		go func(numTries int) {
			count := 0
			for i := 0; i < numTries; i++ {
				x := rand.Float64()
				y := rand.Float64()
				if x*x+y*y <= 1 {
					count++
				}

			}
			countChan <- count
		}(chunk)
		numThread++

		leftIterations -= chunk
		if leftIterations <= 0 {
			break
		}
	}

	totalHits := 0
	for i := 0; i < numThread; i++ {
		totalHits += <-countChan
	}
	close(countChan)
	fmt.Println(4 * float64(totalHits) / float64(totalIterations))
}
