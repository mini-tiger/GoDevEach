package test

import (
	"testing"
	"unsafe"
)

func uuA() string {
	var x []byte = []byte("this is sete1")

	return string(x) // todo  复制内存
}

func uuB() string {
	var x []byte = []byte("this is sete2")
	p := (*string)(unsafe.Pointer(&x)) // todo 共享内存
	//x[0]='t'
	return *p
}

func Benchmark_uu3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuA()
	}
}

func Benchmark_uu4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuB()
	}
}
