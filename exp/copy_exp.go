package main

import "fmt"

func main() {
	src := make([]string, 0)
	src = append(src, "abc")
	src = append(src, "abcd")
	src = append(src, "abcde")
	src = append(src, "abcdef")
	src = append(src, "abcdefg")
	src = append(src, "abcdefgg")
	dst := make([]string, len(src))
	n := copy(dst, src)
	fmt.Printf("copy len = %d\n", n)
	fmt.Printf("dst = %v", dst)
}
