package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"redis-streams/g"
	"strings"
	"time"
)

func consumer(client *redis.Client) {
	// Create a consumer group
	// 消费组中每个g.ConsumerName 名字唯一

	_, err := client.XGroupCreateMkStream(context.Background(), g.StreamName, g.StreamGroup, "0").Result()
	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", g.StreamGroup)
			panic(err)
		}
	}

	go consumerEvent(client)

}

func consumerEvent(client *redis.Client) {
	for {

		// Read a message from the stream
		messages, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    g.StreamGroup,
			Consumer: g.ConsumerName, //组中不同 消费者名字 必须唯一
			Streams:  []string{g.StreamName, ">"},
			/*
				> 符号表示从该 Streams 的最新（最高） ID 开始读取消息
				0：表示读取从 Streams 的最旧（最低） ID 开始读取消息。
				$：表示读取从 Streams 的最高（最新） ID 开始读取消息。
				字符串形式的整数：表示从指定的 ID 开始读取消息。
			*/
			Block: time.Second,
			//Block: 0,     // 如果没有可读数据，XReadGroup 操作将阻塞的时间
			NoAck: false, //可选参数推荐false，true表示不要将消息加入到PEL队列（Pending等待队列）中，相当于在消息读取时直接进行消息确认
			Count: 10,
		}).Result()

		if err != nil {
			fmt.Println("Consumer:", err)
		} else {
			fmt.Println("message len:", len(messages[0].Messages))
			time.Sleep(1 * time.Second)
			fmt.Printf("待处理Pending evnet:%v\n", pendingEventNum(client))
			for _, msg := range messages[0].Messages {
				//fmt.Println(time.Now(), "原始msg:", msg.Values["data"])
				var info g.Msg
				if _, ok := msg.Values["data"]; !ok {
					ackMsg(client, msg)
					fmt.Println("格式不对skip 并ack")
					continue
				}
				err := json.Unmarshal([]byte(msg.Values["data"].(string)), &info)
				ackMsg(client, msg)
				if err != nil {
					fmt.Println("Consumer:", err)
					continue
				}

				fmt.Printf("%v msg info:%v\n", time.Now(), info)
				// Acknowledge the message xxx 确认message

			}
		}
	}
}
func ackMsg(client *redis.Client, msg redis.XMessage) {
	_, err := client.XAck(context.Background(), g.StreamName, g.StreamGroup, msg.ID).Result()
	if err != nil {
		fmt.Println("Consumer:", err)
	}

}

func pendingEventNum(client *redis.Client) int {

	pendingStreams, err := client.XPendingExt(context.Background(), &redis.XPendingExtArgs{
		Stream: g.StreamName,
		Group:  g.StreamGroup,
		Start:  "-",
		End:    "+",
		/*
			-：表示最早的消息。
			+：表示当前时间。
			数字：表示 Unix 时间戳，
		*/
		Count: 100000000000,

		//Consumer string
	}).Result()

	if err != nil {
		panic(err)
	}

	return len(pendingStreams)
	//var streamsRetry []string
	//for _, stream := range pendingStreams {
	//	streamsRetry = append(streamsRetry, stream.ID)
	//}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "172.22.50.25:6379",
		Password: "cc",
		DB:       1,
	})
	consumer(client)
	select {}
}
