package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"

	"PubSub/proto"
	"encoding/json"
	"github.com/micro/go-plugins/broker/nats"
	"log"
)

func main() {

	natsReg :=nats.NewBroker(broker.Addrs("127.0.0.1:4222"))
	//实例化
	service := micro.NewService(
		micro.Name("go.micro.srv"),
		micro.Version("v1.0.0"),
		micro.Broker(natsReg),
	)

	service.Init()

	brok := service.Server().Options().Broker
	if err := brok.Connect(); err != nil {
		log.Fatal(" broker connection failed, error : ", err.Error())
	}

	student := &message.Student{Name: "davie", Classes: "软件工程专业", Grade: 80, Phone: "123456789011"}
	msgBody, err := json.Marshal(student)
	if err != nil {
		log.Fatal(err.Error())
	}

	msg := &broker.Message{
		Header: map[string]string{
			"name": student.Name,
		},
		Body: msgBody,
	}

	err = brok.Publish("go.micro.srv.message", msg)
	if err != nil {
		log.Fatal(" 消息发布失败：%s\n", err.Error())
	} else {
		log.Print("消息发布成功")
	}

	defer brok.Disconnect()
}
