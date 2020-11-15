package main

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func main() {
	prefix := flag.String("prefix", "", "key prefix")
	rev := flag.Int64("rev", 0, "revision of key")
	flag.Parse()

	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := cli.Get(context.TODO(),
		*prefix,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortDescend),
		clientv3.WithMaxCreateRev(*rev),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("head revision = %d\n", resp.Header.Revision)
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %d\n", string(ev.Key), ev.CreateRevision)
	}

}
