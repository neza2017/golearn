package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {
	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		tresp, err := cli.TimeToLive(context.TODO(), resp.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("TTL = %d\n", tresp.TTL)
		time.Sleep(time.Second)
	}
}
