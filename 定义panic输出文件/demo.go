package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

const panicFile = "panic.log"

var globalFile *os.File

func InitPanicFile() error {
	log.Println("init panic file in unix mode")
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	globalFile = file
	if err != nil {
		println(err)
		return err
	}
	fmt.Println(int(file.Fd()), int(os.Stdout.Fd()))
	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}
	if err = syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd())); err != nil {
		return err
	}
	return nil
}

func init() {
	err := InitPanicFile()
	if err != nil {
		println(err)
	}
}

func testPanic() {
	fmt.Println("stdout test") // xxx stdout
	panic("test panic")        // xxx stderr
}

func main() {
	testPanic()
}
