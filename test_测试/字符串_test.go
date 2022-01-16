package test

import (
	"bytes"
	"strings"
	"testing"
)
import "unsafe"

// xxx go test -v -bench=_ -benchmem a_test.go
//func Test_ByteString(t *testing.T) {
//	var x = []byte("Hello World!")
//	var y = *(*string)(unsafe.Pointer(&x))
//	var z = string(x)
//
//	if y != z {
//		t.Fail()
//	}
//}

func Benchmark_Normal(b *testing.B) {
	var x = []byte("Hello")
	x = append(x, []byte(" World!")...)
	for i := 0; i < b.N; i++ {
		_ = string(x)
	}
}

func Benchmark_ByteString(b *testing.B) {
	var x = []byte("Hello")
	x = append(x, []byte(" World!")...)
	for i := 0; i < b.N; i++ {
		_ = *(*string)(unsafe.Pointer(&x))
	}
}

func Benchmark_StringsBuild(b *testing.B) {
	var x strings.Builder
	x.WriteString("Hello ")
	x.WriteString(" World!")

	for i := 0; i < b.N; i++ {
		_ = x.String()
	}
}

func Benchmark_StringsBuildbyte(b *testing.B) {
	var x bytes.Buffer
	x.Write([]byte("Hello "))
	x.Write([]byte("World!"))

	for i := 0; i < b.N; i++ {
		_ = x.String()
	}
}
