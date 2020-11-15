package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func main() {
	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	_, err = cli.Txn(context.TODO()).If().Then(
		clientv3.OpPut("k30", "v30"),
		clientv3.OpPut("k40", "v40"),
	).Commit()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := cli.Txn(context.TODO()).If().Then(
		clientv3.OpGet("k30"),
		clientv3.OpGet("k40"),
	).Commit()
	if err != nil {
		log.Fatal(err)
	}
	for _, rp := range resp.Responses {
		for _, ev := range rp.GetResponseRange().Kvs {
			fmt.Printf("%s -> %s, create revision = %d\n",
				string(ev.Key),
				string(ev.Value),
				ev.CreateRevision)
		}
	}
}
