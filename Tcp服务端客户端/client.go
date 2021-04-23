package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var connFail = make(chan struct{}, 0)
var MiddleChan = make(chan struct{}, 0)

func main() {

	// 接收中间信号    传递重启信号
	go func() {
		for {
			select {
			case <-MiddleChan:
				connFail <- struct{}{}
			}
		}

	}()

	go func() {
		for {
			select {
			case <-connFail:
				log.Println("开始链接")
				StartRecv()
			}
			time.Sleep(1 * time.Second)
		}

	}()
	// 第一次启动
	connFail <- struct{}{}
	select {}

}
func StartRecv() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("链接失败", r.(string))
			MiddleChan <- struct{}{}
		}
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(fmt.Sprintf("dial failed:%s\n", err))
		//os.Exit(1)
	}
	defer conn.Close()
	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			panic(fmt.Sprintf("Read failed:%s\n", err))

		}
		log.Println("count:", n, "msg:", string(buffer))
	}

}
