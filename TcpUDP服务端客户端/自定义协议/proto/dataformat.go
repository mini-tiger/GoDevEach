package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

var Head_Size int = 4
var formatType = binary.BigEndian

// Encode 将消息编码
func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)

	// 写入消息头,整个消息的长度
	err := binary.Write(pkg, formatType, length) //消息头包括 消息的长度 没有包括自身的4字节
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, formatType, []byte(message)) // xxx littleEndian 小头
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	// Peek 返回缓存的一个切片，该切片引用缓存中前 n 字节数据
	// 该操作不会将数据读出，只是引用
	lengthByte, _ := reader.Peek(Head_Size) // 读取前4个字节的数据,也就是读取消息长度
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, formatType, &length) //
	if err != nil {
		return "", err
	}
	fmt.Println(length, reader.Buffered())
	// Buffered返回缓冲中现有的可读取的字节数。
	if int32(reader.Buffered()) < length+4 { // 这里小于代表 数据不够 ，包头中定义的长度
		return "", err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(int32(Head_Size)+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[Head_Size:]), nil // xxx 不显示消息头
}
