package test

import (
	"fmt"
	"runtime"
	"test_ce/utils"
	"testing"
)

// go test -v -bench=Slice -benchmem
// 下面两个要分开使用goland 运行 查看已分配的对象   有道笔记

func Benchmark_Slice1(b *testing.B) {
	b.ResetTimer()
	slice1(b.N)
}

func Benchmark_Slice2(b *testing.B) {
	b.ResetTimer()
	slice2(b.N)
}

func slice1(n int) (ss []int) {
	for i := 0; i < n; i++ {
		s := utils.GenSlice()
		copy(ss, s[2:4])
	}
	//var ss []int
	//copy(ss, s[2:4])
	runtime.GC()
	fmt.Printf("With  (%T), GC took %s\n", ss, utils.TimeGC())
	return
}

func slice2(n int) (ss []int) {
	for i := 0; i < n; i++ {
		s := utils.GenSlice()
		ss = s[2:4]
	}
	//var ss []int
	//copy(ss, s[2:4])
	runtime.GC()
	fmt.Printf("With  (%T), GC took %s\n", ss, utils.TimeGC())
	return
}
