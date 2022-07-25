package initRedis

import "github.com/go-redis/redis"

type RedisClusterObj struct {
	Addr []string
	Auth string
	Db   *redis.ClusterClient
}

// 结构体方法
func (r *RedisClusterObj) initClusterClient() (err error) {
	r.Db = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})

	_, err = r.Db.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
