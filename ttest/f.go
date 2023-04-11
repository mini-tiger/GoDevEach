package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	var str = "ZHONGGUOnihao123"
	strbytes := []byte(str)
	encoded := base64.StdEncoding.EncodeToString(strbytes)

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	decodestr := string(decoded)
	fmt.Println(decodestr, err)
}
