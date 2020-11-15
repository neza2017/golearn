package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func main() {
	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	keyPrefix := "/txn"
	for i := 0; i < 10; i++ {
		_, err := cli.Put(context.TODO(), fmt.Sprintf("%s/key%d", keyPrefix, i), fmt.Sprintf("val%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	resp, err := cli.Get(context.TODO(), keyPrefix, clientv3.WithFirstCreate()...)
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s -> %s\n", string(ev.Key), string(ev.Value))
	}

	resp, err = cli.Get(context.TODO(),
		keyPrefix,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend),
		clientv3.WithLimit(1))
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s -> %s\n", string(ev.Key), string(ev.Value))
	}

}
