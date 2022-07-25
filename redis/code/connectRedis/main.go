package main

import "code/initRedis"

func main() {
	// 实例化RedisSingleObj结构体
	conn := &initRedis.RedisSingleObj{
		Host: "127.0.0.1",
		Port: 6379,
		Auth: "",
	}

	// 初始化连接 Single Redis 服务端
	err := conn.InitSingleRedis()
	if err != nil {
		panic(err)
	}

	// 程序执行完毕释放资源
	defer conn.Db.Close()
}
