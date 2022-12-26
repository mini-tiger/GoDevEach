package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"unsafe"
)

func main() {
	a := 16 << 10
	fmt.Println(a)
	b := []byte("this is byte")
	//fmt.Println(string(b)) // 此方法 会复制[]byte
	s := *(*string)(unsafe.Pointer(&b))
	fmt.Println(s)
	s = "  baad  "
	s = strings.Trim(s, " ")
	debug.PrintStack()
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	fmt.Println(s) // 推荐方法，[]byte 改变，而随着变

}
