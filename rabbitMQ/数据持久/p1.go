package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	// 1. 尝试连接RabbitMQ，建立连接
	// 该连接抽象了套接字连接，并为我们处理协议版本协商和认证等。
	conn, err := amqp.Dial("amqp://user:password@172.16.71.17:5672/")
	if err != nil {
		fmt.Printf("connect to RabbitMQ failed, err:%v\n", err)
		return
	}
	defer conn.Close()

	// 2. 接下来，我们创建一个通道，大多数API都是用过该通道操作的。
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("open a channel failed, err:%v\n", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
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

	// 3. 要发送，我们必须声明要发送到的队列。
	_, err = ch.QueueDeclare(
		"task_queue1", // name
		true,          // 持久的
		false,         // delete when unused
		false,         // 独有的
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		fmt.Printf("declare a queue failed, err:%v\n", err)
		return
	}

	// 4. 然后我们可以将消息发布到声明的队列
	for i := 0; ; i++ {
		time.Sleep(1 * time.Second)
		body := time.Now().Format(time.RFC3339)
		err = ch.Publish(
			"test_queue_ex_topic", // exchange
			"huxiaobai.one111",    // routing key
			false,                 // 立即
			false,                 // 强制
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // 持久
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		if err != nil {
			fmt.Printf("publish a message failed, err:%v\n", err)
			return
		}
		log.Printf(" [x] Sent %s", body)
	}
}

// bodyFrom 从命令行获取将要发送的消息内容
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
