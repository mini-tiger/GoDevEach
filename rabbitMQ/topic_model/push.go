//topic主题模式生产者
package main

import (
	"fmt"
	"rbtmq/rbtmqcs"
	"strconv"
	"time"
)

// https://learnku.com/docs/go-szgbf/1.0/topic-mode-of-rabbitmq-working-mode/8748

func main() {
	rabbitmqOne := rbtmqcs.NewRabbitMQTopic("hxbExc", "huxiaobai.one")
	rabbitmqTwo := rbtmqcs.NewRabbitMQTopic("hxbExc", "huxiaobai.two.cs")
	for i := 0; ; i++ {
		rabbitmqOne.PublishTopic("hello huxiaobai one" + strconv.Itoa(i))
		rabbitmqTwo.PublishTopic("hello huxiaobai two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
