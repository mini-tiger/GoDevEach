package main

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"os"
	"time"
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

	if rdb.Ping(ctx).Err() != nil {
		os.Exit(1)
	}
	str := `{
	"userid": "9678", 
	"username": "杨变利",
	"towerId": "10535",
	"MID": "QTZ63",
	"projectID": "6941",
	"paibanId":"11876",
	"data":[
      {
        "isalarm": "预警",
        "data": 30,
		"fengji": "1",
        "type": "风速",
		"yujing": 30,
        "gaojing": 50,
        "danwei": "m/s"
      },
      {
        "isalarm": "正常",
        "data": 30,
		"action": "上行",
		"actiondata": "10",
        "type": "高度",
		"yujing": 30,
        "gaojing": 50,
        "danwei": "m"
      },
	  
      {
        "isalarm": "告警",
        "data": 30,
		"action": "前进",
		"actiondata": "10",
        "type": "幅度",
		"yujing": 30,
        "gaojing": 50,
        "danwei": "m"
      },
      {
        "isalarm": "告警",
        "data": 30,
        "type": "重量",
		"yujing": 30,
        "gaojing": 50,
        "danwei": "kg"
      },
      {
        "isalarm": "告警",
        "type": "角度",
		"action": "左转",
		"actiondata": "10",
        "data": 30,
		"yujing": 180,
        "gaojing": 30,
        "danwei": "°"
      }
    ]
}`
	go func() {
		for {
			err := rdb.Publish(ctx, "tajiRealTimeData", str).Err()
			if err != nil {
				panic(err)
			}
			//ch := ps.Channel()
			//for msg := range ch {
			//	fmt.Println(msg.Channel, msg.Payload)
			//}
			time.Sleep(500 * time.Millisecond)
		}

	}()

	select {}
}
