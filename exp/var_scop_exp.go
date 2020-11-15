package main

import (
	"fmt"
	"math/rand"
)

func randInt() (int, int) {
	return rand.Int(), rand.Int()
}

func main() {
	p1, i := randInt()
	fmt.Printf("1: &i = %x, %d\n",&i,p1)

	p2, i := randInt()
	fmt.Printf("2: &i = %x, %d\n",&i,p2)
	{
		p3, i := randInt()
		fmt.Printf("3: &i = %x, %d\n",&i,p3)
	}
}
