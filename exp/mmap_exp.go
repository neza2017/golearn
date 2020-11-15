package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	fd, err := os.OpenFile("/tmp/mmap.db", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	const pageSize = 1024

	ref, err := syscall.Mmap(int(fd.Fd()), 0, pageSize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Munmap(ref)

	page := make([]byte, pageSize)
	for i := 0; i < pageSize; i++ {
		page[i] = byte(i % 10)
	}
	fd.WriteAt(page, 0)
	syscall.Fsync((int)(fd.Fd()))
	for i := 0; i < pageSize; i++ {
		if ref[i] != byte(i%10) {
			log.Panicf("no equal at %d \n", i)
		}
	}
	fmt.Println("success")
}
