package main

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func main() {
	prefix := flag.String("prefix", "k", "prefix of key")
	rev := flag.Int64("rev", 0, "revision of key")
	flag.Parse()

	fmt.Printf("prefix=%s,rev=%d\n", *prefix, *rev)

	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	getOps := append(clientv3.WithLastCreate(), clientv3.WithMaxCreateRev(*rev))
	//getOps := append(clientv3.WithLastCreate(), clientv3.WithRev(*rev))
	//getOps := clientv3.WithLastCreate()
	resp, err := cli.Get(context.TODO(), *prefix, getOps...)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Kvs) != 0 {
		fmt.Printf("%s -> %s, rev =%d\n", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value), resp.Header.Revision)
	}

}
