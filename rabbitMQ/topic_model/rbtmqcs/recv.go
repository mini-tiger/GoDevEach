package rbtmqcs

import (
	"fmt"
	"log"
)

func (r *RabbitMQ) RecieveTopic(dc bool) {
	//1.尝试创建交换机exchange 如果交换机存在就不用管他，如果不存在则会创建交换机
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//类型 topic主题模式下我们需要将类型设置为topic
		"topic",
		//进入的消息是否持久化 进入队列如果不消费那么消息就在队列里面 如果重启服务器那么这个消息就没啦 通常设置为false
		dc,
		//是否为自动删除  这里解释的会更加清楚：https://blog.csdn.net/weixin_30646315/article/details/96224842?utm_medium=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromMachineLearnPai2-1.nonecase&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromMachineLearnPai2-1.nonecase
		false,
		//true表示这个exchange不可以被客户端用来推送消息，仅仅是用来进行exchange和exchange之间的绑定
		false,
		//队列消费是否阻塞 fase表示是阻塞 true表示是不阻塞
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha "+"nge")
	//2.试探性创建队列，这里注意队列名称不要写哦
	q, err := r.channel.QueueDeclare(
		//随机生产队列名称 这个地方一定要留空
		r.QueueName,
		dc,
		!dc,
		//具有排他性   排他性的理解 这篇文章还是比较好的：https://www.jianshu.com/p/94d6d5d98c3d
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare a queue")
	//3.绑定队列到exchange中去
	fmt.Printf("%#v\n", q)
	fmt.Printf("%#v\n", r)
	err = r.channel.QueueBind(
		r.QueueName, //队列的名称  通过key去找绑定好的队列
		//在路由模式下，这里的key要填写
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	//4.消费代码
	//4.1接收队列消息
	message, err := r.channel.Consume(
		//队列名称
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答 意思就是收到一个消息已经被消费者消费完了是否主动告诉rabbitmq服务器我已经消费完了你可以去删除这个消息啦 默认是true
		false,
		//是否具有排他性
		false,
		//如果设置为true表示不能将同一个connection中发送的消息传递给同个connectio中的消费者
		false,
		//队列消费是否阻塞 fase表示是阻塞 true表示是不阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//4.2真正开始消费消息
	forever := make(chan bool)
	go func() {
		for d := range message {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
		}
	}()
	fmt.Println("退出请按 ctrl+c")
	<-forever
}
