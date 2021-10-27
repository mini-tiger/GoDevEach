package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"strconv"
	"sync"
	"time"
)

//全局注释 kafka 虽然性能比 rabbitmq要快 但是他丢失数据库的可能性更大，而且还会存在重复接受消息的情况

var Topic = "test-topic"
var partition = int32(0)

func main() {
	sarama.Logger = log{}
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_2_0_0
	cfg.Producer.Return.Errors = true
	cfg.Net.SASL.Enable = false
	cfg.Producer.Return.Successes = true //这个是关键，否则读取不到消息
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewManualPartitioner //允许指定分组
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Offsets.AutoCommit = struct {
		Enable   bool
		Interval time.Duration
	}{Enable: false, Interval: 1 * time.Second}

	//cfg.Group.Return.Notifications = true
	cfg.ClientID = "service-exchange-api"
	var kafka = KafkaConfig{
		Addrs:  []string{"192.168.43.111:9092"},
		Config: cfg,
	}
	_, _, err := NewKafkaClient(kafka)
	fmt.Println("err:", err)
}

//发送消息 此为异步发送消息
func NewAsyncProducer(client sarama.Client, i int) error {
	c, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return err
	}

	defer c.Close()
	p, o, err := c.SendMessage(&sarama.ProducerMessage{Topic: Topic, Value: sarama.StringEncoder("消息发送成功拉ssssssss！！！！" + strconv.Itoa(i))})
	if err != nil {
		fmt.Printf("err:", err)
		return err
	}
	fmt.Printf("partition:%v, offset:%v\n", p, o)
	/*c, err := sarama.NewAsyncProducerFromClient(client)
	//sarama.NewSyncProducerFromClient() 此为同步
	if err != nil {
		return  err
	}
	defer  c.Close()
	//Topic 为主题，Partition为区域 Partition如果不给默认为0 记得设置cfg.Producer.Partitioner = sarama.NewManualPartitioner 这里为允许设置指定的分区
	//分区是从0开始，记得在启动配置文件时修改Partition的分区
	//不同的主题包括不同的分区都是有着不同的offset
	c.Input() <-  &sarama.ProducerMessage{Topic: Topic,Key:sarama.StringEncoder(fmt.Sprintf("/topic/market/order-trade")), Value: sarama.StringEncoder("消息发送成功拉ssssssss！！！！"+strconv.Itoa(i))}
	select {
	//case msg := <-producer.Successes():
	//    log.Printf("Produced message successes: [%s]\n",msg.Value)
	case err := <-c.Errors():
		fmt.Println("Produced message failure: ", err)
	default:
		//fmt.Println("Produced message success",err)
	}*/
	return nil
}

//客户端接收消息
func NewKafkaClient(cfg KafkaConfig) (sarama.Client, func(), error) {

	//创建链接 创建客户机
	c, err := sarama.NewClient(cfg.Addrs, cfg.Config)
	if err != nil {
		return nil, nil, err
	}
	offsetid, _ := c.GetOffset("test-topic", 0, sarama.OffsetNewest)
	fmt.Printf("topic: test-topic, 最新offset:%v\n", offsetid)

	go func() {
		//目前默认是肯定能使用的
		consumer, err := sarama.NewConsumerGroupFromClient("default-group", c)
		//client, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "group-1", cfg.Config)
		if err != nil {
			fmt.Println(err)
		}
		loopConsumer(consumer, Topic, partition, "b")
		consumer.Close()
	}()

	//go func() {
	//	for i := 0; i < 10; i++ {
	//		NewAsyncProducer(c, i)
	//	}
	//}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return c, func() {
		err := c.Close()
		if err != nil {
			fmt.Print(err)
		}
	}, nil
}

func loopConsumer(consumer sarama.ConsumerGroup, topic string, partition int32, item string) {
	go func() {
		for err := range consumer.Errors() {
			fmt.Println(err)
		}
	}()
	ctx, _ := context.WithCancel(context.Background())
	hand := MainHandler{}
	for {
		err := consumer.Consume(ctx, []string{"test-topic"}, &hand)

		if err != nil {
			fmt.Println(err)
			break
		}
		if ctx.Err() != nil {
			break
		}
	}
	/*for {
		msg := <-partitionConsumer.Messages()
		pom.MarkOffset(msg.Offset + 1, "备注")
		fmt.Printf("[%s] : Consumed message: [%s], offset: [%d]\n",item, string(msg.Value), msg.Offset)
	}*/

}

type KafkaConfig struct {
	Addrs  []string
	Config *sarama.Config
}

type MainHandler struct {
}

func (m *MainHandler) Setup(sess sarama.ConsumerGroupSession) error {
	// 如果极端情况下markOffset失败，需要手动同步offset
	return nil
}

func (m *MainHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	// 如果极端情况下markOffset失败，需要手动同步offset
	return nil
}

//此方法会自动控制偏移值，当分组里的主题消息被接收到时，则偏移值会进行加1 他是跟着主题走的
func (m *MainHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s Offset  = %v\n", string(message.Value), message.Timestamp, message.Topic, message.Offset)
		sess.MarkMessage(message, "")

		if message.Offset%3 == 0 {
			// 手动提交，不能频繁调用，耗时9ms左右，macOS i7 16GB
			//t1 := time.Now().Nanosecond()
			sess.Commit()
			//t2 := time.Now().Nanosecond()
			//fmt.Println("commit cost:", (t2-t1)/(1000*1000), "ms")
		}
	}
	return nil
}

type log struct{}

func (log) Print(v ...interface{}) {
	fmt.Println(v...)
}

func (log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (log) Println(v ...interface{}) {
	fmt.Println(v...)
}
