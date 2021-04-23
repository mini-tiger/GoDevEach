package main

import (
	"context"
	"encoding/json"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"time"
)

// xxx https://github.com/go-redis/redis

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.40.111:6379",
		Password: "Root1q2w", // no password set
		DB:       0,          // use default DB
	})
}

func StrExampleClient() {

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

type A struct {
	D string
}

func listExampleClient() {
	data, _ := json.Marshal(A{"1"})
	if err := rdb.RPush(ctx, "queue", data).Err(); err != nil {
		panic(err)
	}

	// use `rdb.BLPop(0, "queue")` for infinite waiting time
	result, err := rdb.BLPop(ctx, 1*time.Second, "queue").Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("redis key:%s,Value:%+v\n", result[0], result[1])
}

func main() {
	StrExampleClient()
	listExampleClient()
}
