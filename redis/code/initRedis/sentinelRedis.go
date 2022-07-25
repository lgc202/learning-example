package initRedis

import "github.com/go-redis/redis"

type RedisSentinelObj struct {
	MasterName    string
	SentinelAddrs []string
	Auth          string
	Db            *redis.Client
}

// 结构体方法
func (r *RedisSentinelObj) initSentinelClient() (err error) {
	r.Db = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    r.MasterName,
		SentinelAddrs: []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"},
		Password:      r.Auth,
	})
	_, err = r.Db.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
