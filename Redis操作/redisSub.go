package main

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"os"
	"time"

	"fmt"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.43.177:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

}

func main() {
	//client1, err := redis.CreateClient(8, "192.168.43.11:6379", "")
	//if err != nil {
	//	log.Println("NewRedis err:", err)
	//}
	//defer client1.Close()
	//Redis111.Conn = client1
	if rdb.Ping(ctx).Err() != nil {
		os.Exit(1)
	}

	go func() {
		ps := rdb.Subscribe(ctx, "tajiRealTimeData")
		_, err := ps.Receive(ctx)
		ch := ps.Channel()
		if err != nil {
			panic(err)
		}
		for {

			for msg := range ch {
				fmt.Printf("%s : %s\n", time.Now(), msg.Payload)
			}
		}

	}()

	select {}
}
