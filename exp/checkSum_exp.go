package main

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"unsafe"
)

type metaStruct struct {
	d0       int64
	d1       int64
	d2       int64
	checkSum uint64
}

func (m *metaStruct) sum64() uint64 {
	h := fnv.New64a()
	p2 := (*[unsafe.Offsetof(metaStruct{}.checkSum)]byte)(unsafe.Pointer(m))
	_, _ = h.Write(p2[:])
	fmt.Printf("type of p2 is %v \n", reflect.ValueOf(p2).Kind())
	fmt.Printf("offset of checksum = %d\n", unsafe.Offsetof(metaStruct{}.checkSum))
	fmt.Printf("size of metaStruct = %d\n", unsafe.Sizeof(metaStruct{}))
	return h.Sum64()
}

func main() {
	var ms metaStruct
	ms.d0 = -1
	ms.d1 = -1
	ms.d2 = -1
	ms.checkSum = ms.sum64()
	fmt.Println(ms)
}
