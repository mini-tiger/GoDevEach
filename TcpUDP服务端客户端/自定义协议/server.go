package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"自定义协议/proto"
)

func process(conn net.Conn) {
	defer func() {
		log.Printf("client Addr:%s , close", conn.RemoteAddr())
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	for {
		msg, err := proto.Decode(reader) // xxx 消息解码
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", msg)
	}
}

func main() {

	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	log.Println("start recv")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}
