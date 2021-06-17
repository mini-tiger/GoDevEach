package main

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const FILENAME1 = "test.mmap"

func main() {
	pagesize := os.Getpagesize()

	file, _ := os.OpenFile(FILENAME1, os.O_RDWR|os.O_CREATE, 0644)

	state, _ := file.Stat()

	if state.Size() == 0 {
		_, _ = file.WriteAt(bytes.Repeat([]byte{'0'}, pagesize), 0)

		state, _ = file.Stat()
	}

	fmt.Printf("pagesize: %d\n filesize: %d\n", pagesize, state.Size())

	data, _ := unix.Mmap(int(file.Fd()), 0, int(state.Size()), unix.PROT_WRITE, unix.MAP_SHARED)

	fmt.Printf("文件内容:%v\n", string(data))
}
