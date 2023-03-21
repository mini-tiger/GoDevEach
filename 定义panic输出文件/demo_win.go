package main

import (
	"log"
	"os"
	"syscall"
)

const (
	kernel32dll = "C:\\Windows\\System32\\kernel32.dll"
)

const panicFile = "C:/panic.log"

var globalFile *os.File

func InitPanicFile() error {
	log.Println("init panic file in windows mode")
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND, 0666)
	globalFile = file
	if err != nil {
		return err
	}
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	v, _, err := setStdHandle.Call(uintptr(sh), uintptr(file.Fd()))
	if v == 0 {
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
	panic("test panic")
}

func main() {
	testPanic()
}
