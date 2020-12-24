package main

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int, charset string) string {
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	threads := 1000
	iterations := 1000000
	completeChan := make(chan int)
	fmt.Println("********begin")
	for i := 0; i < threads; i++ {
		go func() {
			for k := 0; k < iterations; k++ {
				time.Sleep(time.Nanosecond)
				//fmt.Println(randomString(k, charset))
			}
			completeChan <- 1
		}()
	}

	for i := 0; i < threads; i++ {
		<-completeChan
	}

	fmt.Println("********end")
}
