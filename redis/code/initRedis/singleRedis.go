package initRedis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisSingleObj struct {
	Host     string
	Port     uint16
	Auth     string
	Database int
	Db       *redis.Client
}

func (r *RedisSingleObj) InitSingleRedis() (err error) {
	// Redis连接格式拼接
	redisAddr := fmt.Sprintf("%s:%d", r.Host, r.Port)
	// Redis 连接对象: NewClient将客户端返回到由选项指定的Redis服务器。
	r.Db = redis.NewClient(&redis.Options{
		Addr:        redisAddr,  // redis服务ip:port
		Password:    r.Auth,     // redis的认证密码
		DB:          r.Database, // 连接的database库
		IdleTimeout: 300,        // 默认Idle超时时间
		PoolSize:    100,        // 连接池
	})
	fmt.Printf("Connecting Redis : %v\n", redisAddr)

	// 验证是否连接到redis服务端
	res, err := r.Db.Ping().Result()
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return err
	} else {
		fmt.Printf("Connect Successful! Ping => %v\n", res)
		return nil
	}
}
