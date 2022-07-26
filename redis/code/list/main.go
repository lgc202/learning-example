package main

import (
	"code/initRedisV8"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func LPushExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LPushExample============================")
	rdb.FlushDB(ctx) // 清空数据库

	// 插入指定值到list列表中，返回值是当前列表元素的数量
	// 使用LPush()方法将数据从左侧压入链表（后进先出）,也可以从右侧压如链表对应的方法是RPush()
	count, _ := rdb.LPush(ctx, "list", 1, 2, 3).Result()
	fmt.Println("插入到list集合中元素的数量: ", count)
	fmt.Println(rdb.LRange(ctx, "list", 0, -1))
}

func LInsertExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LInsertExample============================")
	rdb.FlushDB(ctx) // 清空数据库

	// LInsert() 在某个位置插入新元素
	rdb.LPush(ctx, "list", 1, 2, 3).Result() // [3 2 1]
	// 在名为key的缓存项值为2的元素前面插入一个值，值为123 ， 注意只会执行一次
	_ = rdb.LInsert(ctx, "list", "before", "2", 123).Err()
	fmt.Println(rdb.LRange(ctx, "list", 0, -1))
	// 在名为key的缓存项值为2的元素后面插入一个值，值为321
	_ = rdb.LInsert(ctx, "list", "after", "2", 321).Err()
	fmt.Println(rdb.LRange(ctx, "list", 0, -1))
}

func LSetExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LSetExample============================")
	rdb.FlushDB(ctx)                         // 清空数据库
	rdb.LPush(ctx, "list", 1, 2, 3).Result() // [3 2 1]

	// LSet() 设置某个元素的值,下标是从0开始的
	val1, _ := rdb.LSet(ctx, "list", 2, 256).Result()
	fmt.Println("是否成功将下标为2的元素值改成256: ", val1)
	fmt.Println(rdb.LRange(ctx, "list", 0, -1))
}

func LLenExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LLenExample============================")
	rdb.FlushDB(ctx)                         // 清空数据库
	rdb.LPush(ctx, "list", 1, 2, 3).Result() // [3 2 1]

	// LLen() 获取链表元素个数
	length, _ := rdb.LLen(ctx, "list").Result()
	fmt.Printf("当前链表的长度为: %v\n", length)
}

func LIndexExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LIndexExample============================")
	rdb.FlushDB(ctx)                         // 清空数据库
	rdb.LPush(ctx, "list", 1, 2, 3).Result() // [3 2 1]

	// LIndex() 获取链表下标对应的元素
	val2, _ := rdb.LIndex(ctx, "list", 2).Result()
	fmt.Printf("下标为2的值为: %v\n", val2)
}

func LPopExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LPopExample============================")
	rdb.FlushDB(ctx)                         // 清空数据库
	rdb.LPush(ctx, "list", 1, 2, 3).Result() // [3 2 1]

	// 从链表左侧弹出数据
	val3, _ := rdb.LPop(ctx, "list").Result()
	fmt.Printf("弹出下标为0的值为: %v\n", val3)
}

func LRemExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start LPopExample============================")
	rdb.FlushDB(ctx)                            // 清空数据库
	rdb.LPush(ctx, "list", 1, 2, 3, 3).Result() // [3 3 2 1]

	// LRem() 根据值移除元素 lrem key count value
	n, _ := rdb.LRem(ctx, "list", 2, "3").Result()
	fmt.Printf("移除了: %v 个\n", n)
	fmt.Println(rdb.LRange(ctx, "list", 0, -1))
}

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
	fmt.Println()

	ctx := context.Background()
	LPushExample(redisClient, ctx)
	LInsertExample(redisClient, ctx)
	LSetExample(redisClient, ctx)
	LLenExample(redisClient, ctx)
	LIndexExample(redisClient, ctx)
	LPopExample(redisClient, ctx)
	LRemExample(redisClient, ctx)
}
