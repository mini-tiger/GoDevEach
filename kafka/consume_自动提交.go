package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup
var topic = "kafkaMsgTypeOs"

// xxx 不能 从上次偏移量 继续
func main() {
	sarama.Logger = log.New(os.Stdout, "sarama: ", log.LstdFlags)

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0
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
	}{Enable: true, Interval: 3 * time.Second} // 3秒 提交1次

	//创建新的消费者
	//consumer, err := sarama.NewConsumer([]string{"192.168.43.111:9092"}, config)

	client, err := sarama.NewClient([]string{"172.16.71.17:9092"}, cfg)
	if err != nil {
		fmt.Println("fail to start consumer", err)
	}

	offsetid, _ := client.GetOffset("kafkaMsgTypeOs", 0, sarama.OffsetNewest)
	//
	fmt.Printf("最新 offset:%v\n", offsetid)
	consumer, err := sarama.NewConsumerFromClient(client)

	//根据topic获取所有的分区列表
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("fail to get list of partition,err:", err)
	}
	fmt.Println("partitionList:", partitionList)
	//遍历所有的分区
	for p := range partitionList {
		//针对每一个分区创建一个对应分区的消费者
		pc, err := consumer.ConsumePartition(topic, int32(p), sarama.OffsetNewest) // xxx offsetOldset 全部数据
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", p, err)
		}
		defer pc.AsyncClose()
		wg.Add(1)
		//异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("partition:%d Offset:%d Key:%v Value:%s \n",
					msg.Partition, msg.Offset, msg.Key, msg.Value)

			}
		}(pc)
	}
	wg.Wait()
}
