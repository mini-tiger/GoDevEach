package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

// xxx https://github.com/go-redis/redis

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "172.22.50.25:6379",
		Password: "cc", // no password set
		DB:       1,    // use default DB
	})
}

var migrateParser *viper.Viper
var rkey string = "a:b:c"

func newViperParserFromFile(target string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(path.Base(target))
	v.AddConfigPath(path.Dir(target))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(v); err != nil {
			fmt.Println(err)
		}

		fmt.Println(v.AllSettings())
		//err = rdb.Set(ctx, rkey, data, 0).Err()
		//if err != nil {
		//	panic(err)
		//}
	})
	return v, nil
}
func SetMigrateFromFile(target string) (*viper.Viper, error) {
	var err error

	// /data/migrate.yaml -> /data/migrate
	split := strings.Split(target, ".")
	filePath := split[0]
	migrateParser, err = newViperParserFromFile(filePath)
	if err != nil {

		return migrateParser, err
	}
	return migrateParser, nil
}

func StrToRedisExampleClient() {
	file := "/data/work/go/GoDevEach/Redis操作/common.yaml"

	cc, _ := SetMigrateFromFile(file)
	fmt.Printf("webServer.api:%v\n", cc.Get("webServer.api"))

	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(content))

	err = rdb.Set(ctx, rkey, string(content), 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, rkey).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
	time.Sleep(100 * time.Second)

	//val2, err := rdb.Get(ctx, "key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
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
	StrToRedisExampleClient()
	listExampleClient()
}
