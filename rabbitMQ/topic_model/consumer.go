package main

import "rbtmq/rbtmqcs"

func main() {
	//#号表示匹配多个单词 也就是读取hxbExc交换机里面所有队列的消息
	rabbitmq := rbtmqcs.NewRabbitMQTopic("hxbExc", "#")
	rabbitmq.RecieveTopic()

	//这里只是匹配到了huxiaobai.后边只能是一个单词的key 通过这个key去找绑定到交换机上的相应的队列
	//rabbitmq := rbtmqcs.NewRabbitMQTopic("hxbExc","huxiaobai.*.cs")
	//rabbitmq.RecieveTopic()

}
