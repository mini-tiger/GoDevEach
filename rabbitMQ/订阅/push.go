package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// 连接rabbitmq
	conn, err := amqp.Dial("amqp://user:password@172.16.71.17:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建信道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明交换机
	err = ch.ExchangeDeclare(
		"tizi3651", // 交换机名字
		"fanout",   // 交换机类型，fanout发布订阅模式
		true,       // 是否持久化
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// xxx
	_, err = ch.QueueDeclare(
		"haha1", // name
		true,    // 持久的
		false,   // delete when unused
		false,   // 独有的
		false,   // no-wait
		nil,     // arguments
	)
	// 消息内容

	for {
		body := fmt.Sprintf("Hello Tizi365.com time:%s!", time.Now().Format(time.RFC3339))
		// 推送消息
		err = ch.Publish(
			"tizi3651", // exchange（交换机名字，跟前面声明对应）
			"",         // 路由参数，fanout类型交换机，自动忽略路由参数，填了也没用。
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain", // 消息内容类型，这里是普通文本
				Body:         []byte(body), // 消息内容
			})

		log.Printf("发送内容 %s ", body)
		time.Sleep(200 * time.Millisecond)
	}

}
