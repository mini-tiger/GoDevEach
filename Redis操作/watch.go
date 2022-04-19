package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "172.22.50.25:6379",
		Password: "cc", // no password set
		DB:       0,    // use default DB
	})
}

// xxx redis乐观锁支持，可以通过watch监听一些Key, 如果这些key的值没有被其他人改变的话，才可以提交事务。
// xxx 不能监听CRUD
func main() {
	fn := func(tx *redis.Tx) error {
		// 先查询下当前watch监听的key的值
		v, err := tx.Get(ctx, "key").Result()
		if err != nil && err != redis.Nil {
			return err
		}

		// 这里可以处理业务
		fmt.Println(v)

		// 如果key的值没有改变的话，Pipelined函数才会调用成功
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在这里给key设置最新值
			pipe.Set(ctx, "key", "new value", 0)
			return nil
		})
		return err
	}

	// 使用Watch监听一些Key, 同时绑定一个回调函数fn, 监听Key后的逻辑写在fn这个回调函数里面
	// 如果想监听多个key，可以这么写：client.Watch(fn, "key1", "key2", "key3")
	_ = rdb.Watch(ctx, fn, "key")
}
