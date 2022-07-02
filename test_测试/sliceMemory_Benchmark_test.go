package test

import (
	"testing"
)

// go test -v -bench=Slice -benchmem
// 下面两个要分开使用goland 运行 查看已分配的对象   有道笔记

func genSlice() []int {
	s := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		s[i] = i
	}
	return s
}

func Benchmark_Slice1(b *testing.B) {
	s := genSlice()
	for i := 0; i < b.N; i++ {
		_ = slice1(s)
	}
}

func Benchmark_Slice2(b *testing.B) {
	s := genSlice()
	for i := 0; i < b.N; i++ {
		_ = slice2(s)
	}
}

func slice1(s []int) (ss []int) {

	//var ss []int
	copy(ss, s[2:4])
	return
}

func slice2(s []int) []int {

	return s[2:4]
}
