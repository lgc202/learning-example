package initRedisV8

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
)

type RedisSingleObj struct {
	Host     string
	Port     uint16
	Auth     string
	Database int
}

func (r *RedisSingleObj) InitSingleRedis() (*redis.Client, error) {
	// Redis连接格式拼接
	redisAddr := fmt.Sprintf("%s:%d", r.Host, r.Port)
	// Redis 连接对象: NewClient将客户端返回到由选项指定的Redis服务器。
	RedisClient = redis.NewClient(&redis.Options{
		Addr:        redisAddr,  // redis服务ip:port
		Password:    r.Auth,     // redis的认证密码
		DB:          r.Database, // 连接的database库
		IdleTimeout: 300,        // 默认Idle超时时间
		PoolSize:    100,        // 连接池
	})
	fmt.Printf("Connecting Redis : %v\n", redisAddr)

	// go-redis库v8版本相关命令都需要传递context.Context参数,Background 返回一个非空的Context,它永远不会被取消，没有值，也没有期限。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 验证是否连接到redis服务端
	res, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return nil, err
	}

	// 输出连接成功标识
	fmt.Printf("Connect Successful! \nPing => %v\n", res)
	return RedisClient, nil
}
