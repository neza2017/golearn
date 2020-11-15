package main

import (
	"log"
	"sort"
)

type myStructExp struct {
	id int
}
type myStructExps []*myStructExp

func (s myStructExps) Len() int {
	return len(s)
}

func (s myStructExps) Less(i, j int) bool {
	return s[i].id < s[j].id
}

func (s myStructExps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	var ss myStructExps
	for i := 10; i > 0; i-- {
		d := &myStructExp{id: i}
		ss = append(ss, d)
	}
	for _, s := range ss {
		log.Println(s.id)
	}
	log.Println("==========")
	sort.Sort(ss)
	for _, s := range ss {
		log.Println(s.id)
	}

}
