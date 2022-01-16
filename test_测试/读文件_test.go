package test

import (
	"io/ioutil"
	"os"
	"testing"
	"unsafe"
)

/**
 * @Author: Tao Jun
 * @Description: test
 * @File:  读文件_test
 * @Version: 1.0.0
 * @Date: 2021/6/17 下午5:25
 */
const FILENAME1 = "/home/go/GoDevEach/mmap共享内存/test.mmap"

// xxx go test -v -bench=ReadFile -benchmem 读文件_test.go

//func ReadFile() []byte {
//	file, _ := os.OpenFile(FILENAME1, os.O_RDWR|os.O_CREATE, 0644)
//
//	state, _ := file.Stat()
//	data, _ := unix.Mmap(int(file.Fd()), 0, int(state.Size()), unix.PROT_READ, unix.MAP_SHARED) // xxx 不能重复运行
//	return data
//	//return
//}
//
//func Benchmark_ReadFile(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = ReadFile()
//	}
//}

func ReadFile1() string {
	//file, _ := os.OpenFile(FILENAME1, os.O_RDWR|os.O_CREATE, 0644)

	data2, _ := ioutil.ReadFile(FILENAME1)
	//data1:=make([]byte,1024)
	//file.Read(data1)
	//fmt.Println(string(data1))
	//fmt.Println(string(data2))

	return *(*string)(unsafe.Pointer(&data2))
}
func Benchmark_ReadFile1(b *testing.B) {
	//pagesize := os.Getpagesize()
	for i := 0; i < b.N; i++ {

		_ = ReadFile1()

	}
}

func ReadFile2() string {
	file, _ := os.OpenFile(FILENAME1, os.O_RDWR|os.O_CREATE, 0644)

	//data2,_:=ioutil.ReadFile(FILENAME1)
	data1 := make([]byte, 1024)
	file.Read(data1)
	//fmt.Println(string(data1))
	//fmt.Println(string(data2))

	return *(*string)(unsafe.Pointer(&data1))
}
func Benchmark_ReadFile2(b *testing.B) {
	//pagesize := os.Getpagesize()
	for i := 0; i < b.N; i++ {

		_ = ReadFile2()

	}
}
