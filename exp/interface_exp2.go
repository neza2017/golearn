package main

import "fmt"

type baseA interface {
	Print()
}

type structA struct {
	val int
}

func (s *structA) Print() {
	fmt.Printf("struct A val = %d\n", s.val)
}

type structB struct {
	a structA
}

func doPrint(b baseA) {
	b.Print()
}

func main() {
	a := structA{val: 10}
	doPrint(&a)
	//b := structB{a: structA{20}}
	//doPrint(&b)
}
