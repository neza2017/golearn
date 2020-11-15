package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"sync"
)

type myMutex struct {
	client  *clientv3.Client
	leaseId clientv3.LeaseID
	pfx     string
	myKey   string
	myRev   int64
}

func newMyMutex(cli *clientv3.Client, pfx string) (*myMutex, error) {
	resp, err := cli.Grant(context.TODO(), 15)
	if err != nil {
		return nil, err
	}
	if _, err := cli.KeepAlive(context.TODO(), resp.ID); err != nil {
		return nil, err
	}

	mu := &myMutex{
		client:  cli,
		leaseId: resp.ID,
		pfx:     pfx,
		myKey:   fmt.Sprintf("%s%x", pfx, resp.ID),
		myRev:   -1,
	}

	return mu, nil
}

func (mu *myMutex) Close() error {
	ctx, cancel := context.WithCancel(context.TODO())
	_, err := mu.client.Revoke(ctx, mu.leaseId)
	cancel()
	return err
}

func (mu *myMutex) UnLock(ctx context.Context) error {
	_, err := mu.client.Delete(ctx, mu.myKey)
	return err
}

func (mu *myMutex) Lock(ctx context.Context) error {
	cmp := clientv3.Compare(clientv3.CreateRevision(mu.myKey), "=", 0)
	put := clientv3.OpPut(mu.myKey, "", clientv3.WithLease(mu.leaseId))
	get := clientv3.OpGet(mu.myKey)
	getOwner := clientv3.OpGet(mu.pfx, clientv3.WithFirstCreate()...)
	resp, err := mu.client.Txn(ctx).If(cmp).Then(put, getOwner).Else(get, getOwner).Commit()
	if err != nil {
		return err
	}
	mu.myRev = resp.Header.Revision
	if !resp.Succeeded {
		mu.myRev = resp.Responses[0].GetResponseRange().Kvs[0].CreateRevision
	}
	ownerRev := resp.Responses[1].GetResponseRange().Kvs[0].CreateRevision
	if ownerRev == mu.myRev {
		return nil
	}
	getOps := append(clientv3.WithLastCreate(), clientv3.WithMaxCreateRev(mu.myRev-1))
	for {
		resp, err := mu.client.Get(ctx, mu.pfx, getOps...)
		if err != nil {
			return err
		}
		if len(resp.Kvs) == 0 {
			break
		}
		wch := mu.client.Watch(ctx, string(resp.Kvs[0].Key), clientv3.WithRev(resp.Header.Revision))
		for wr := range wch {
			isBreak := false
			for _, ev := range wr.Events {
				if ev.Type == mvccpb.DELETE {
					isBreak = true
					break
				}
			}
			if isBreak {
				break
			}
		}
	}
	return nil
}

//func main() {
//	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//	mu, err := newMyMutex(cli, "/my_mutex")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer mu.Close()
//
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		i := i
//		go func() {
//			defer wg.Done()
//			if err := mu.Lock(context.TODO()); err != nil {
//				log.Fatal(err)
//			}
//			fmt.Printf("goroutine %d locked\n", i)
//		}()
//	}
//	wg.Wait()
//	fmt.Println("exit main")
//}

func main() {
	var cnt int64
	var wg sync.WaitGroup
	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
			if err != nil {
				log.Fatal(err)
			}
			defer cli.Close()
			mu, err := newMyMutex(cli, "/my_mutex/")
			if err != nil {
				log.Fatal(err)
			}
			defer mu.Close()
			mu.Lock(context.TODO())

			defer mu.UnLock(context.TODO())

			for k := 0; k < 100000; k++ {
				cnt = cnt + 1
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(cnt)
}

//func main() {
//	endpoints := []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//
//	m1, err := NewMyMutex(cli, "/my_mutex/")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer m1.Close()
//
//	if err = m1.Lock(context.TODO()); err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("lock by m1")
//	m1.UnLock(context.TODO())
//
//	m2, err := NewMyMutex(cli, "/my_mutex/")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer m2.Close()
//	if err = m2.Lock(context.TODO()); err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("finish")
//
//}
