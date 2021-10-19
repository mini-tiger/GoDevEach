package main

import (
	"TCP_optimization/proto"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  client
 * @Version: 1.0.0
 * @Date: 2021/4/25 上午11:32
 */
func main() {
	data := []byte("[这里才是一个完整的数据包]")
	l := len(data)
	fmt.Println("data length: ", l)
	magicNum := make([]byte, proto.SplitLenMark)          // 分隔符标识长度
	binary.BigEndian.PutUint32(magicNum, proto.SplitMark) // 分隔符

	lenNum := make([]byte, proto.DataLenMark)     //数据标识长度
	binary.BigEndian.PutUint16(lenNum, uint16(l)) //数据内容长度

	packetBuf := bytes.NewBuffer(magicNum)
	packetBuf.Write(lenNum)
	packetBuf.Write(data)

	conn, err := net.DialTimeout("tcp", "localhost:4044", time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for {
		_, err = conn.Write(packetBuf.Bytes())
		if err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}
	}
}
