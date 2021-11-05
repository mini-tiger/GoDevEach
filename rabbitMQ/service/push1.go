package service

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("push %s: %s", msg, err)
	}
}

func PublishExchageQueue(conn *amqp.Connection, exchange string, routingKeys []string, queueName string) {
	// 1. 建立RabbitMQ连接
	//conn, err := amqp.Dial("amqp://user:password@172.16.71.17:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2. 创建channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3. 声明exchange,routing key
	//exchange := "test_topic_exchange"
	//routingKey1 := "user.save"
	//routingKey2 := "user.update"
	//routingKey3 := "user.delete.abc"

	// 4. 声明（创建）一个交换机
	//name:交换器的名称。
	//kind:也叫作type，表示交换器的类型。有四种常用类型：direct、fanout、topic、headers。
	//durable:是否持久化，true表示是。持久化表示会把交换器的配置存盘，当RMQ Server重启后，会自动加载交换器。
	//autoDelete:是否自动删除，true表示是。至少有一条绑定才可以触发自动删除，当所有绑定都与交换器解绑后，会自动删除此交换器。
	//internal:是否为内部，true表示是。客户端无法直接发送msg到内部交换器，只有交换器可以发送msg到内部交换器。
	//noWait:是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args:直接写nil，没研究过，不解释。

	//注意，在生产者里声不声明(创建)交换机都可以。这里声明，是为了防止消费者没有启动或者这个交换机原先不存在，导致消息投递丢失。

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	msg := "Hello World RabbitMQ Topic Exchange Message ... "

	// 5. 发送消息
	//exchange：要发送到的交换机名称，对应图中exchangeName。
	//key：路由键，对应图中RoutingKey。
	//mandatory：消息发布的时候设置消息的 mandatory 属性用于设置消息在发送到交换器之后无法路由到队列的情况对消息的处理方式， 设置为 true 表示将消息返回到生产者，否则直接丢弃消息。直接false，不建议使用。
	//immediate ：参数告诉服务器至少将该消息路由到一个队列中，否则将消息返回给生产者。immediate参数告诉服务器，如果该消息关联的队列上有消费者，则立刻投递:如果所有匹配的队列上都没有消费者，则直接将消息返还给生产者，不用将消息存入队列而等待消费者了。直接false，不建议使用。RabbitMQ 3.0版本开始去掉了对immediate参数的支持。
	//msg：要发送的消息，msg对应一个Publishing结构，Publishing结构里面有很多参数，这里只强调几个参数，其他参数暂时列出，但不解释。
	for index, routingKey := range routingKeys {
		err = ch.Publish(
			exchange,   // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("%s %d", msg, index)),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf("[info] [x] Sent %s %d", msg, index)
	}

}
