package main

import (
	"fmt"
	"math/rand"
)

func randInt2() (int, int) {
	return rand.Int(), rand.Int()
}

func main() {
	v1, v2 := randInt2()
	fmt.Printf("%x : %x\n", &v1, &v2)
	if v1, v3 := randInt2(); v1 != 0 {
		fmt.Printf("%x : %x\n", &v1, &v3)
	}
	if v1, v2 := randInt2(); v1 != 0 {
		fmt.Printf("%x : %x\n", &v1, &v2)
	}
	if v1, v2 = randInt2(); v1 != 0 {
		fmt.Printf("%x : %x\n", &v1, &v2)
	}
}
