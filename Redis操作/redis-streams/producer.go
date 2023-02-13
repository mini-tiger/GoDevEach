package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"redis-streams/g"
	"time"
)

func producer(client *redis.Client) {
	var i uint64
	for i = 1; ; i++ {
		// Write a message to the stream
		stream := client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: g.StreamName,
			//Values: map[string]interface{}{
			//	"message": fmt.Sprintf("message %d", i),
			//},
			Values: map[string]interface{}{
				"data": createMsg(i),
			},
		})

		fmt.Println("Producer:", stream.Val(), stream.String(), stream.Err())
		time.Sleep(1 * time.Second)
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "172.22.50.25:6379",
		Password: "cc",
		DB:       1,
	})
	producer(client)
}

func createMsg(num uint64) []byte {
	b, err := json.Marshal(&g.Msg{Val: fmt.Sprintf("message %d", num), ID: num})
	if err != nil {
		return []byte{}
	}
	return b
}
