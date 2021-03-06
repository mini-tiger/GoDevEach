package rbtmqcs

import "github.com/streadway/amqp"

//topic主题模式step1:创建RabbitMQ实例
func NewRabbitMQTopic(exchange string, routingkey string, queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	return NewRabbitMQ(queueName, exchange, routingkey)
}

//topic主题模式step2:发送消息
func (r *RabbitMQ) PublishTopic(message string, dc bool) {
	//1.尝试创建交换机
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
	if dc {
		_, err = r.channel.QueueDeclare(
			r.QueueName, // name
			dc,          // 持久的
			!dc,         // delete when unused
			false,       // 独有的
			false,       // no-wait
			nil,         // arguments
		)
	}
	//2.发送消息
	ppp := amqp.Transient
	if dc {
		ppp = amqp.Persistent
	}
	err = r.channel.Publish(
		r.Exchange,
		//除了设置交换机这也要设置绑定的key值
		r.Key,
		//如果为true 会根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返还给发送者
		false,
		//如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者则会把消息返还给发送者
		false,
		amqp.Publishing{
			DeliveryMode: ppp,
			ContentType:  "text/plain",
			Body:         []byte(message), //发送的内容一定要转换成字节的形式
		})
}
