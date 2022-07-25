package main

import (
	"code/initRedisV8"
	"context"
	"fmt"
)

// 最新版本的go-redis库的相关命令都需要传递context.Context参数
var ctx = context.Background()

func main() {
	// 实例化RedisSingleObj结构体
	conn := &initRedisV8.RedisSingleObj{
		Host: "127.0.0.1",
		Port: 6379,
		Auth: "",
	}

	// 初始化连接 Single Redis 服务端
	redisClient, err := conn.InitSingleRedis()
	if err != nil {
		fmt.Printf("[Error] - %v\n", err)
		return
	}

	// 程序执行完毕释放资源
	defer redisClient.Close()
}
