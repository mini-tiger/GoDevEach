package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://user:password@172.16.71.17:5672/")
	if err != nil {
		fmt.Printf("connect to RabbitMQ failed, err:%v\n", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("open a channel failed, err:%v\n", err)
		return
	}
	defer ch.Close()

	_ = ch.ExchangeDeclare(
		//交换机名称
		"test_queue_ex_topic",
		//类型 topic主题模式下我们需要将类型设置为topic
		"topic",
		//进入的消息是否持久化 进入队列如果不消费那么消息就在队列里面 如果重启服务器那么这个消息就没啦 通常设置为false
		true,
		//是否为自动删除  这里解释的会更加清楚：https://blog.csdn.net/weixin_30646315/article/details/96224842?utm_medium=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromMachineLearnPai2-1.nonecase&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromMachineLearnPai2-1.nonecase
		false,
		//true表示这个exchange不可以被客户端用来推送消息，仅仅是用来进行exchange和exchange之间的绑定
		false,
		//队列消费是否阻塞 fase表示是阻塞 true表示是不阻塞
		false,
		nil,
	)

	// 声明一个queue
	q, err := ch.QueueDeclare(
		"task_queue1", // name
		true,          // 声明为持久队列
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Printf("ch.Qos() failed, err:%v\n", err)
		return
	}

	// 立即返回一个Delivery的通道
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // 注意这里传false,关闭自动消息确认
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Printf("ch.Consume failed, err:%v\n", err)
		return
	}

	// 开启循环不断地消费消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//dotCount := bytes.Count(d.Body, []byte("."))
			//t := time.Duration(dotCount)
			//time.Sleep(t * time.Second)
			//log.Printf("Done")
			d.Ack(false) // 手动传递消息确认
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
