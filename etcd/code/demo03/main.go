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
	if _, err = kv.Put(context.TODO(), "/cron/jobs/job1", "job1", clientv3.WithPrevKV()); err != nil {
		panic(err)
	}

	// 写入另外一个Job
	if _, err = kv.Put(context.TODO(), "/cron/jobs/job2", "job2", clientv3.WithPrevKV()); err != nil {
		panic(err)
	}

	// 读取/cron/jobs/为前缀的所有key
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		panic(err)
	}

	// 获取成功, 我们遍历所有的kvs
	fmt.Println(getResp.Kvs)
}
