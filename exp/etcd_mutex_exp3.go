package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"log"
	"sync"
	"time"
)

func main() {
	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	muPrefix := "/mutext/test/"
	ss, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer ss.Close()
	mu := concurrency.NewMutex(ss, muPrefix)
	err = mu.Lock(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	mu2 := concurrency.NewMutex(ss, muPrefix)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			err := mu2.Lock(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("goroutine %d lock %s\n", i, muPrefix)
			wg.Done()
		}()
	}
	time.Sleep(time.Second)
	fmt.Printf("main UnLock %s\n", muPrefix)
	mu.Unlock(context.TODO())
	wg.Wait()
	fmt.Println("exit main")

}
