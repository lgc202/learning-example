package main

import (
	"code/initRedisV8"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// Redis String Set/Get 示例
func setGetExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start setGetExample============================")
	// 1.Set 设置 key 如果设置为0则表示永不过期
	rdb.FlushDB(ctx) // 清空数据库
	err := rdb.Set(ctx, "setGetExample", 100, 60*time.Second).Err()
	if err != nil {
		fmt.Printf("set setGetExample failed, err:%v\n", err)
		panic(err)
	}

	// 2.Get 获取已存在的Key其存储的值
	val, err := rdb.Get(ctx, "setGetExample").Result() // 获取其值
	if err != nil {
		fmt.Printf("get setGetExample failed, err:%v\n", err)
		panic(err)
	}
	fmt.Printf("val -> setGetExample ：%v\n", val)
}

func nilExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start nilExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	_, err := rdb.Get(ctx, "nilExample").Result()
	if err == redis.Nil {
		fmt.Println("[ERROR] - Key [nilExample] not exist")
	} else if err != nil {
		fmt.Printf("get nilExample failed, err:%v\n", err)
		panic(err)
	}
}

func existsExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start existsExample============================")
	// Exists() 方法用于检测某个key是否存在
	rdb.FlushDB(ctx) // 清空数据库
	n, _ := rdb.Exists(ctx, "existsExample").Result()
	if n > 0 {
		fmt.Println("existsExample key 存在!")
	} else {
		fmt.Println("existsExample key 不存在!")
		rdb.Set(ctx, "existsExample", "weiyi", 60*time.Second)
	}

	n, _ = rdb.Exists(ctx, "existsExample").Result()
	if n > 0 {
		fmt.Println("existsExample key 存在!")
	}
}

func setNxExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start setNxExample============================")
	// SetNX 当不存在key时将进行设置该可以并设置其过期时间
	rdb.FlushDB(ctx) // 清空数据库
	val, err := rdb.SetNX(ctx, "setNxExample", "helloworld", 0).Result()
	if err != nil {
		fmt.Printf("set username failed, err:%v\n", err)
		panic(err)
	}
	fmt.Printf("set setNxExample success: %v\n", val) // true

	val, err = rdb.SetNX(ctx, "setNxExample", "helloworld", 0).Result()
	if err != nil {
		fmt.Printf("set setNxExample failed, err:%v\n", err)
		panic(err)
	}
	fmt.Printf("set setNxExample success: %v\n", val) // false
}

func keysExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start keysExample============================")
	// Keys() 根据正则获取keys
	rdb.FlushDB(ctx) // 清空数据库
	rdb.Set(ctx, "k1", "v1", 60*time.Second)
	rdb.Set(ctx, "k2", "v2", 60*time.Second)

	keys, _ := rdb.Keys(ctx, "*").Result()
	fmt.Printf("All Keys: %v \n", keys)

	// 根据前缀获取Key
	keys, _ = rdb.Keys(ctx, "k*").Result()
	fmt.Printf("All Keys: %v \n", keys)
}

func dbSizeExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start dbSizeExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// DBSize() 查看当前数据库key的数量
	rdb.Set(ctx, "k1", "v1", 60*time.Second)
	rdb.Set(ctx, "k2", "v2", 60*time.Second)
	num, _ := rdb.DBSize(ctx).Result()
	fmt.Printf("Keys number : %v \n", num)
}

func typeExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start typeExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// Type() 方法用户获取一个key对应值的类型
	rdb.Set(ctx, "username", "v1", 60*time.Second)
	vType, err := rdb.Type(ctx, "username").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("username key type : %v\n", vType)
}

func expireExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start expireExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// Expire()方法是设置某个时间段(time.Duration)后过期，ExpireAt()方法是在某个时间点(time.Time)过期失效
	rdb.Set(ctx, "k1", "v1", 0)
	rdb.Set(ctx, "k2", "v2", 0)
	val4, _ := rdb.Expire(ctx, "k1", time.Minute*2).Result()
	if val4 {
		fmt.Println("k1 过期时间设置成功", val4)
	} else {
		fmt.Println("k1 过期时间设置失败", val4)
	}

	val5, _ := rdb.ExpireAt(ctx, "k2", time.Now().Add(time.Minute*2)).Result()
	if val5 {
		fmt.Println("k2 过期时间设置成功", val5)
	} else {
		fmt.Println("k2 过期时间设置失败", val5)
	}
}

func ttlExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start ttlExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// TTL()与 PTTL()方法可以获取某个键的剩余有效期, 区别是精确度不一样
	rdb.Set(ctx, "k1", "v1", 2*time.Minute)
	rdb.Set(ctx, "k2", "v2", 2*time.Minute)

	k1TTL, _ := rdb.TTL(ctx, "k1").Result() // 获取其key的过期时间
	k2TTL, _ := rdb.PTTL(ctx, "k2").Result()
	fmt.Printf("k1 TTL : %v, k2 TTL : %v\n", k1TTL, k2TTL)
}

func delExample(rdb *redis.Client, ctx context.Context) {
	// keys命令优点： 花的时间短
	// keys命令的缺点：这个命令会阻塞redis多路复用的io主线程，如果这个线程阻塞，在此执行之间其他的发送向redis服务端的命令，都会阻塞，从而引发一系列级联反应，导致瞬间响应卡顿，从而引发超时等问题，所以应该在生产环境禁止用使用keys和类似的命令smembers，这种时间复杂度为O（N），且会阻塞主线程的命令，是非常危险的。
	//
	// scan 命令优点：不阻塞 IO主线程
	// scan命令缺点：
	// 1.由于scan采用的增量迭代，当redis中的key是随时变化的，比如key增加减少或者key的名字变更，这种情况，scan就暴露他的弊端了，可能无法获取所有的key了
	// 2. 返回的数据有可能重复
	// 所以说增量式迭代命令只能提供有限的保证！

	fmt.Println("start delExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// 当通配符匹配的key的数量不多时，可以使用Keys()得到所有的key在使用Del命令删除。
	rdb.Set(ctx, "k1", "v1", 60*time.Second)
	rdb.Set(ctx, "k2", "v2", 60*time.Second)

	keys, _ := rdb.Keys(ctx, "*").Result()
	num, err := rdb.Del(ctx, keys...).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Del() : ", num)

	// 如果key的数量非常多的时候，我们可以搭配使用Scan命令和Del命令完成删除
	rdb.Set(ctx, "k1", "v1", 60*time.Second)
	rdb.Set(ctx, "k2", "v2", 60*time.Second)
	iter := rdb.Scan(ctx, 0, "k*", 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}

func flushDBExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start flushDBExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	rdb.Set(ctx, "k1", "v1", 60*time.Second)
	rdb.Set(ctx, "k2", "v2", 60*time.Second)

	// 清空当前数据库，因为连接的是索引为0的数据库，所以清空的就是0号数据库
	flag, err := rdb.FlushDB(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("FlushDB() : ", flag)
}

func getRangeExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start getRangeExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	err := rdb.Set(ctx, "hello", "Hello World!", 0).Err()
	if err != nil {
		panic(err)
	}

	// GetRange ：字符串截取, 截取的是value
	// 注：即使key不存在，调用GetRange()也不会报错，只是返回的截取结果是空"",可以使用fmt.Printf("%q\n", val)来打印测试
	val1, _ := rdb.GetRange(ctx, "hello", 1, 7).Result()
	fmt.Printf("key: hello, value: %v\n", val1) //截取到的内容为:  ello Wo
}

func appendExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start appendExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// Append()表示往字符串后面追加元素，返回值是字符串的总长度
	rdb.Set(ctx, "hello", "Hello World!", 0)
	length, _ := rdb.Append(ctx, "hello", " Go Programer").Result()
	val, _ := rdb.Get(ctx, "hello").Result()
	fmt.Printf("当前缓存key的长度为: %v，值: %v \n", length, val)
}

// stringIntExample 数据类型演示
func stringIntExample(rdb *redis.Client, ctx context.Context) {
	fmt.Println("start stringIntExample============================")
	rdb.FlushDB(ctx) // 清空数据库
	// 设置整形的key
	err := rdb.SetNX(ctx, "number", 1, 0).Err()
	if err != nil {
		panic(err)
	}
	// Incr()、IncrBy()都是操作数字，对数字进行增加的操作
	// Decr()、DecrBy()方法是对数字进行减的操作，和Incr正好相反
	// incr是执行原子加1操作
	val3, _ := rdb.Incr(ctx, "number").Result()
	fmt.Printf("Incr -> key当前的值为: %v\n", val3) // 2
	// incrBy是增加指定的数
	val4, _ := rdb.IncrBy(ctx, "number", 6).Result()
	fmt.Printf("IncrBy -> key当前的值为: %v\n", val4) // 8

	// StrLen 也可以返回缓存key的对应value长度
	length2, _ := rdb.StrLen(ctx, "number").Result()
	fmt.Printf("number 值长度: %v\n", length2)
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

	// String 数据类型操作
	ctx := context.Background()
	setGetExample(redisClient, ctx)
	nilExample(redisClient, ctx)
	existsExample(redisClient, ctx)
	setNxExample(redisClient, ctx)
	keysExample(redisClient, ctx)
	dbSizeExample(redisClient, ctx)
	typeExample(redisClient, ctx)
	expireExample(redisClient, ctx)
	ttlExample(redisClient, ctx)
	delExample(redisClient, ctx)
	flushDBExample(redisClient, ctx)
	getRangeExample(redisClient, ctx)
	appendExample(redisClient, ctx)
	stringIntExample(redisClient, ctx)
}
