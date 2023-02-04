package main

import (
	"fmt"
	"github.com/felipeagger/go-redis-streams/common"
	"math/rand"
	"time"

	evt "github.com/felipeagger/go-redis-streams/packages/event"
	"github.com/felipeagger/go-redis-streams/packages/utils"
	"github.com/go-redis/redis/v7"
)

//xxx https://blog.csdn.net/li_qinging/article/details/117450163
// https://pkg.go.dev/github.com/go-redis/redis/v7
var (
	streamName string = common.StreamName
	client     *redis.Client
)

func init() {
	var err error

	client, err = utils.NewRedisClient(&redis.Options{
		Addr:     "172.22.50.25:6379",
		Password: "cc", // no password set
		DB:       1,    // use default DB
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	generateEvent()
}

func generateEvent() {
	var userID uint64 = 0
	for i := 0; i < 10; i++ {

		userID++ //uint64(rand.Intn(1000))

		eventType := []evt.Type{evt.LikeType, evt.CommentType}[rand.Intn(2)]

		if eventType == evt.LikeType {

			newID, err := produceMsg(map[string]interface{}{
				"type": string(eventType),
				"data": &evt.LikeEvent{
					Base: &evt.Base{
						Type:     eventType,
						DateTime: time.Now(),
					},
					UserID: userID,
				},
			})

			checkError(err, newID, string(eventType), userID)

		} else {

			comment := []string{"Go e Top!", "Go e demais!", "Go e vida!"}[rand.Intn(3)]

			newID, err := produceMsg(map[string]interface{}{
				"type": string(eventType),
				"data": &evt.CommentEvent{
					Base: &evt.Base{
						Type:     eventType,
						DateTime: time.Now(),
					},
					UserID:  userID,
					Comment: comment,
				},
			})

			checkError(err, newID, string(eventType), userID, comment)
		}

	}
}

func produceMsg(event map[string]interface{}) (string, error) {

	return client.XAdd(&redis.XAddArgs{
		Stream: streamName,
		Values: event,
	}).Result()
}

func checkError(err error, newID, eventType string, userID uint64, comment ...string) {
	if err != nil {
		fmt.Printf("produce event error:%v\n", err)
	} else {

		if len(comment) > 0 {
			fmt.Printf("produce event success Type:%v UserID:%v Comment:%v offset:%v\n",
				string(eventType), userID, comment, newID)
		} else {
			fmt.Printf("produce event success Type:%v UserID:%v offset:%v\n",
				string(eventType), userID, newID)
		}

	}
}
