package main

import (
	"context"
	"flag"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func main() {
	cr := flag.Int("createRevision", 0, "create revision")
	flag.Parse()

	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	keyPrefix := "/txn_key"
	// if the keyPrefix is not exist ,then the CreateRevision of this key is 0

	resp, err := cli.Txn(context.TODO()).If(
		clientv3.Compare(clientv3.CreateRevision(keyPrefix), "=", *cr),
	).Then(
		clientv3.OpPut(keyPrefix, "success"),
	).Else(
		clientv3.OpPut(keyPrefix, "failed"),
	).Commit()

	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)

}
