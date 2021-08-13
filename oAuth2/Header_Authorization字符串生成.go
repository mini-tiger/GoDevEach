package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	msg := "Tom:admin@123" //11个字符应该装成15个base64编码的字符
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Printf("Basic %s\n", encoded) //SGVsbG8sd29ybGQ=,后面的=是作填充用的
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(string(decoded)) //Hello,world
}
