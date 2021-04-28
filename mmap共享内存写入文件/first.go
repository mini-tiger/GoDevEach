package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

const FILENAME = "test.mmap"

func main() {
	pagesize := os.Getpagesize()

	file, _ := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE, 0644)

	state, _ := file.Stat()

	if state.Size() == 0 {
		_, _ = file.WriteAt(bytes.Repeat([]byte{'0'}, pagesize), 0)

		state, _ = file.Stat()
	}

	fmt.Printf("pagesize: %d\n filesize: %d\n", pagesize, state.Size())

	data, _ := unix.Mmap(int(file.Fd()), 0, int(state.Size()), unix.PROT_WRITE, unix.MAP_SHARED)

	// 关闭文件，不影响mmap
	file.Close()

	//for i := 0; i < 8; i++ {
	//data[i] = '1'
	//}

	for i, x := range strings.Split("hello world", "") {
		data[i] = []byte(x)[0] // string 转 byte ,用[]byte(*)[index]的方式
		//fmt.Println(i, x)
	}

	unix.Munmap(data)
}
