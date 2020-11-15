package main

import (
	"fmt"
	"reflect"
)

func main() {
	a1 := make([]int, 16)
	a2 := new([16]int)
	a3 := [16]int{}
	a4 := a1[:5]
	a5 := a2[:5]
	a6 := a3[:5]
	//a7 := a1[:20]
	//a8 := a2[:20]
	//a9 := a3[:20]
	fmt.Printf("a1 type = %v\n", reflect.ValueOf(a1).Kind())
	fmt.Printf("a2 type = %v\n", reflect.ValueOf(a2).Kind())
	fmt.Printf("a3 type = %v\n", reflect.ValueOf(a3).Kind())
	fmt.Printf("a4 type = %v\n", reflect.ValueOf(a4).Kind())
	fmt.Printf("a5 type = %v\n", reflect.ValueOf(a5).Kind())
	fmt.Printf("a6 type = %v\n", reflect.ValueOf(a6).Kind())
	b1 := a1[:]
	fmt.Printf("addr %x, addr %x\n", &b1[0], &a1[0])
}
