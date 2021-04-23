package main

import (
	"PubSub/proto"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/nats"

	"log"
)

func main() {
	// xxx 先启动nats   nats --port 4222 --addr 127.0.0.1
	natsReg :=nats.NewBroker(broker.Addrs("127.0.0.1:4222"))
	//实例化
	service := micro.NewService(
		micro.Name("go.micro.srv"),
		micro.Version("v1.0.0"),
		micro.Broker(natsReg),
	)

	//初始化
	service.Init()

	//订阅事件
	pubSub := service.Server().Options().Broker
	if err := pubSub.Connect(); err != nil {
		log.Fatal(" broker connect failed , error is : %v\n", err)
	}

	_, err := pubSub.Subscribe("go.micro.srv.message", func(event broker.Event) error {
		var req *message.StudentRequest
		fmt.Println(string(event.Message().Body))

		if err := json.Unmarshal(event.Message().Body, &req); err != nil {
			return err
		}
		fmt.Println(" 接收到信息：", req)
		//去执行其他操作
		return nil
	})

	if err != nil {
		log.Printf("sub error: %v\n", err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
