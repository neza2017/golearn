package main

import "fmt"

func varArgs(arg ...int) {
	for _, a := range arg {
		fmt.Println(a)
	}
}

func main() {
	varArgs(1, 2, 3, 4, 5)
	fmt.Println("===========")
	varArgs([]int{6, 7, 8, 9, 0}...)
}
