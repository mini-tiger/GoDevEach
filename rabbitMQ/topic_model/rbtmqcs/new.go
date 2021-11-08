package rbtmqcs

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//url 格式 amqp://账号:密码@rabbitmq服务器地址:端口号/Virtual Host
//格式在golang语言当中是固定不变的
const MQURL = "amqp://user:password@172.16.71.17:5672/"

//const MQURL = "amqp://huj:123456@192.168.1.219:5672/rbtmqcs"
type RabbitMQ struct {
	conn    *amqp.Connection //需要引入amqp包 https://learnku.com/articles/44185教会你如何引用amqp包
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//链接信息
	Mqurl string
}

//断开channel和connection
//为什么要断开channel和connection 因为如果不断开他会始终使用和占用我们的channel和connection 断开是为了避免资源浪费
func (r *RabbitMQ) Destory() {
	r.channel.Close() //关闭信道资源
	r.conn.Close()    //关闭链接资源
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建RabbitMQ结构体实例
func NewRabbitMQ(queuename string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queuename, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	//创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)  //通过amqp.Dial()方法去链接rabbitmq服务端
	rabbitmq.failOnErr(err, "创建连接错误!")              //调用我们自定义的failOnErr()方法去处理异常错误信息
	rabbitmq.channel, err = rabbitmq.conn.Channel() //链接上rabbitmq之后通过rabbitmq.conn.Channel()去设置channel信道
	rabbitmq.failOnErr(err, "获取channel失败!")
	return rabbitmq
}
