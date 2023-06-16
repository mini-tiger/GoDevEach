package main

import (
	"configcenter/src/storage/dal/redis"
	"fmt"
	"time"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "172.22.50.25:6379",
		Username:     "neolink",
		Password:     "Ne01ink2022!",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           7,
	})
}

func ExampleNewClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	pong, err := rdb.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleParseURL() {
	opt, err := redis.ParseURL("redis://:qwerty@localhost:6379/1")
	if err != nil {
		panic(err)
	}
	fmt.Println("addr is", opt.Addr)
	fmt.Println("db is", opt.DB)
	fmt.Println("password is", opt.Password)

	// Create client as usually.
	_ = redis.NewClient(opt)

	// Output: addr is localhost:6379
	// db is 1
	// password is qwerty
}

func main() {
	doCommand()
}

// doCommand go-redis基本使用示例
func doCommand() {
	//ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	//defer cancel()

	// 执行命令获取结果
	val, err := rdb.Get("key").Result()
	fmt.Println(val, err)

	// 先获取到命令对象
	cmder := rdb.Get("key")
	fmt.Println(cmder.Val()) // 获取值
	fmt.Println(cmder.Err()) // 获取错误

	// 直接执行命令获取错误
	err = rdb.Set("key", 10, time.Hour).Err()

	// 直接执行命令获取值
	value := rdb.Get("key").Val()
	fmt.Println(value)
}
