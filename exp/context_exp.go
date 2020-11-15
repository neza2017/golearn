package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func simpleContextRoutine(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("start %s\n", name)
	<-ctx.Done()
	fmt.Printf("exit %s\n", name)
}

func main() {
	base, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		ctx := context.WithValue(base, fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
		go simpleContextRoutine(ctx, fmt.Sprintf("routine %d", i), &wg)
	}
	time.Sleep(time.Second)

	cancel()
	fmt.Printf("wait group\n")
	wg.Wait()
	fmt.Printf("exit main\n")
}
