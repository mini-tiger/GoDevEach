package g

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"sync"
)

var RDB *redis.Client
var ctx = context.Background()

var cfglock sync.Mutex

func GetRDB() *redis.Client {
	cfglock.Lock()
	defer cfglock.Unlock()
	return RDB
}

func init() {
	if RDB == nil { // eee 单例模式
		RDB = redis.NewClient(&redis.Options{
			Addr:     "192.168.40.111:6379",
			Password: "Root1q2w",
			DB:       0,
		})
	}

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	pong, err := RDB.Ping(ctx).Result()
	fmt.Println(pong, err)

}
