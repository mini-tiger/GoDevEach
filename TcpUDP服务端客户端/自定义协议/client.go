package main

import (
	"fmt"
	"net"
	"自定义协议/proto"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  client
 * @Version: 1.0.0
 * @Date: 2021/4/25 上午11:32
 */

// socket_stick/client2/main.go

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := fmt.Sprintf(`num:%d,Hello, Hello. How are you?`, i)
		data, err := proto.Encode(msg) // xxx 消息编码
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			return
		}
		conn.Write(data)
	}

}
