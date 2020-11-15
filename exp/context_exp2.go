package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	fmt.Println(time.Now())
	go func() {
		fmt.Printf("start go routine\n")
		<-ctx.Done()
		fmt.Printf("end go routine\n")
	}()
	for i := 0; i < 10; i++ {
		t, ok := ctx.Deadline()
		if ok {
			fmt.Println(t)
		}
		time.Sleep(time.Second)
	}
	fmt.Printf("exit main\n")

}
