package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	defer client.Close()
	fmt.Println("连接etcd成功:", client.Endpoints())

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hello", clientv3.WithPrevKV()); err != nil {
		panic(err)
	}

	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "bye", clientv3.WithPrevKV()); err != nil {
		panic(err)
	}

	fmt.Println("Revision:", putResp.Header.Revision)
	if putResp.PrevKv != nil { // 打印hello
		fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
	}

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1" /*clientv3.WithCountOnly()*/); err != nil {
		panic(err)
	}

	fmt.Printf("Kvs=%v\nCount=%v\n", getResp.Kvs, getResp.Count)
}
