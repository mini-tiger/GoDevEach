//topic主题模式生产者
package main

import (
	"fmt"
	"log"
	"rbtmq/rbtmqcs"
	"strconv"
	"time"
)

// https://learnku.com/docs/go-szgbf/1.0/topic-mode-of-rabbitmq-working-mode/8748

func main() {
	for i := 0; ; i++ {
		begin()
		log.Println("loop :", i)
	}
}
func begin() {
	rabbitmqOne := rbtmqcs.NewRabbitMQTopic("hxbExc12", "huxiaobai.one", "abc1")
	//rabbitmqTwo := rbtmqcs.NewRabbitMQTopic("hxbExc", "huxiaobai.two.cs","abc123")
	for i := 0; i < 10; i++ {
		rabbitmqOne.PublishTopic("hello huxiaobai one"+strconv.Itoa(i), false) //xxx dc  true,消费者 重启 也能 连接上 断开之前的数据， false  如果发送端 不持久化队列，消费端重启，则接收更新的数据
		time.Sleep(200 * time.Millisecond)
		fmt.Println(i)
	}
	rabbitmqOne.Destory()
}
