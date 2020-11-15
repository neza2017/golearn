package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

func randPage(page []byte) {
	size := int(unsafe.Sizeof(uint64(0)))
	uintLen := len(page) / size
	for i := 0; i < uintLen-1; i++ {
		dPtr := (*uint64)(unsafe.Pointer(&page[i*size]))
		*dPtr = rand.Uint64()
	}
}

func hashPage(page []byte) uint64 {
	hashEnd := len(page) - int(unsafe.Sizeof(uint64(0)))
	h := fnv.New64a()
	_, _ = h.Write(page[:hashEnd])
	return h.Sum64()
}

func setPage(page []byte) {
	randPage(page)
	hashEnd := len(page) - int(unsafe.Sizeof(uint64(0)))
	dPtr := (*uint64)(unsafe.Pointer(&page[hashEnd]))
	*dPtr = hashPage(page)

}

func main() {
	pageFile := "/tmp/page.db"
	if _, err := os.Stat(pageFile); err == nil {
		os.Remove(pageFile)
	}

	pageSize := syscall.Getpagesize()
	fmt.Printf("page size = %d \n", pageSize)
	fd, err := os.OpenFile(pageFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	ref, err := syscall.Mmap(int(fd.Fd()), 0, pageSize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Munmap(ref)

	page := make([]byte, pageSize)
	setPage(page)
	if n, err := fd.WriteAt(page, 0); n != pageSize || err != nil {
		log.Fatal(err)
	}
	if err := syscall.Fsync(int(fd.Fd())); err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		hashEnd := pageSize - int(unsafe.Sizeof(uint64(0)))
		dPtr := (*uint64)(unsafe.Pointer(&ref[hashEnd]))
		ret := false
		cnt := 0
		var initHash uint64
		initHash = 0
		for {
			if ret {
				log.Printf("read count = %d\n", cnt)
				return
			}
			h1 := hashPage(ref)
			if initHash == h1 {
				log.Printf("same hash value, i = %x, h = %x\n", initHash, h1)
			}
			initHash = h1
			if *dPtr != h1 {
				log.Fatalf("check sum error, h = %x, d = %x\n", h1, *dPtr)
			}
			select {
			case <-ch:
				ret = true
			default:
				cnt++
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			setPage(page)
			if n, err := fd.WriteAt(page, 0); n != pageSize || err != nil {
				log.Fatal(err)
			}
			if err := syscall.Fsync(int(fd.Fd())); err != nil {
				log.Fatal(err)
			}
		}
		ch <- struct{}{}
	}()
	wg.Wait()

}
