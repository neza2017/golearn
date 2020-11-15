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
	_, err = cli.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 15; i++ {
		gresp, err := cli.Get(context.TODO(), "foo", clientv3.WithPrefix())
		if err != nil {
			log.Fatal(err)
		}
		for _, ev := range gresp.Kvs {
			fmt.Printf("%s -> %s, ", string(ev.Key), string(ev.Value))
		}
		tresp, err := cli.TimeToLive(context.TODO(), resp.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ttl = %d\n", tresp.TTL)
		time.Sleep(time.Second)
		if i == 2 {
			fmt.Printf("keep alive TTL = %d\n", tresp.TTL)
			_, err = cli.KeepAlive(context.TODO(), resp.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("exit main\n")
}
