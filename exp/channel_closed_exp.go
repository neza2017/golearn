package main

import (
	"log"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	go func() {
		for i := 0; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		total := 0
		for {
			if val, ok := <-ch; ok {
				total += val
			} else {
				break
			}
		}
		log.Printf("recv total = %d", total)
	}()
	wg.Wait()
}
