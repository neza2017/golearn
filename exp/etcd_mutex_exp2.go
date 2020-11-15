package main

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"log"
	"time"
)

func main() {
	cliName := flag.String("cliName", "mutexClinet", "client name")
	keyPrefix := flag.String("keyPrefix", "/mutex/lock/", "key prefix")
	safeClose := flag.Bool("safeClose", false, "safeClose")
	flag.Parse()

	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	if *safeClose {
		defer func() {
			fmt.Printf("close etcd client %s\n", *cliName)
			cli.Close()
		}()
	}
	ss, err := concurrency.NewSession(cli, concurrency.WithTTL(1))
	if err != nil {
		log.Fatal(err)
	}
	if *safeClose {
		defer func() {
			fmt.Printf("close etcd session %s\n", *cliName)
		}()
	}

	mu := concurrency.NewMutex(ss, *keyPrefix)

	err = mu.Lock(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("lock mutex on %s\n", *cliName)
	time.Sleep(5 * time.Second)
	fmt.Printf("%s exit\n", *cliName)
	fmt.Println("======================")
}
