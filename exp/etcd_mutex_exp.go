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

const (
	mutexKey = "/mutex/lock"
)

func mutexLock(endpoints []string, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s start to lock\n", name)
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	ss, err := concurrency.NewSession(cli, concurrency.WithTTL(1))
	if err != nil {
		log.Fatal(err)
	}
	defer ss.Close()
	mu := concurrency.NewMutex(ss, mutexKey)

	err = mu.Lock(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("lock mutex on %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("%s exit\n", name)
}

func mutexTryLock(endpoints []string, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	ss, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer ss.Close()
	mu := concurrency.NewMutex(ss, mutexKey)
	for {
		err := mu.TryLock(context.TODO())
		if err == nil {
			break
		}
		if err == concurrency.ErrLocked {
			fmt.Printf("%s try to lock, but failed\n", name)
			time.Sleep(time.Second)
		} else {
			log.Fatal(err)
		}
	}
	fmt.Printf("lock mutex on %s\n", name)
	time.Sleep(5 * time.Second)
	fmt.Printf("%s exit\n", name)
}

func main() {
	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go mutexLock(endpoints, fmt.Sprintf("client_%d", i), &wg)
	}
	wg.Wait()
	time.Sleep(time.Second)
	fmt.Println("==================")
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go mutexTryLock(endpoints, fmt.Sprintf("client_%d", i), &wg)
	}
	wg.Wait()
	time.Sleep(time.Second)
	fmt.Println("exit main")
}
