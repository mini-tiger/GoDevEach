package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"sync"
	"time"
)

func main() {
	// 定义要ping的IP地址列表
	ips := []string{
		"8.8.8.8",
		"8.8.4.4",
		"208.67.222.222",
		"208.67.220.220",
		"172.24.50.25", // not conn
		"223.5.5.5",
	}

	var wg sync.WaitGroup
	wg.Add(len(ips))

	// 用于接收ping结果的通道
	results := make(chan string, len(ips))

	// 遍历IP列表，对每个IP进行ping操作
	for _, ip := range ips {
		go func(ip string) {
			defer wg.Done()
			pinger, err := ping.NewPinger(ip)
			if err != nil {
				fmt.Printf("Error creating pinger: %v\n", err)
				results <- fmt.Sprintf("%s: Error", ip)
				return
			}

			pinger.Count = 3                  // 设置发送的ping包数量
			pinger.Timeout = time.Second      // 设置超时时间
			pinger.Interval = time.Second / 2 // 设置发送间隔

			pinger.OnFinish = func(stats *ping.Statistics) {
				results <- fmt.Sprintf("%s: %v, %v", ip, stats.AvgRtt, stats.PacketLoss)
			}

			pinger.Run()
		}(ip)
	}

	// 等待所有goroutines完成
	wg.Wait()
	close(results)

	// 打印ping结果
	for result := range results {
		fmt.Println(result)
	}
}
