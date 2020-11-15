package main

import "fmt"

func typeSwitchFoo(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("int value = %d\n", v)
	case string:
		fmt.Printf("string value = %s\n", v)
	case *int:
		fmt.Printf("*int, addr = %x, value = %d\n", v, *v)
	case *string:
		fmt.Printf("*string, addr = %x, value = %s\n", v, *v)
	case []int:
		fmt.Printf("[]int, value = ")
		for n := range v {
			fmt.Printf(" %d", n)
		}
		fmt.Printf("\n")
	}

}

func main() {
	x := 5
	s := "hello world"
	typeSwitchFoo(x)
	typeSwitchFoo(s)
	typeSwitchFoo(&x)
	typeSwitchFoo(&s)
	typeSwitchFoo([]int{1, 2, 3, 4, 5, 6})
}
