package test

import (
	"fmt"
	"runtime"
	"test_ce/utils"
	"testing"
)

// 下面两个要分开使用goland 运行 查看已分配的对象   有道笔记

/*
go test -v -bench=Append1 -benchmem -cpu=1 -run=none -count=1
go test -v -bench=Append2 -benchmem -cpu=1 -run=none -count=1
go test -v -bench=Copy -benchmem -cpu=1 -run=none -count=1
*/
func Benchmark_Experiment3Append1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 20; j++ {
			s = append(s, []int{j, j + 1, j + 2, j + 3, j + 4}...)
		}
	}
	runtime.GC()
	fmt.Printf(" GC took %s\n", utils.TimeGC())
}

func Benchmark_Experiment3Append2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 100) // xxx cap 固定,不超过 不会多分配 内存次数 ，可变长
		for j := 0; j < 20; j++ {
			s = append(s, []int{j, j + 1, j + 2, j + 3, j + 4}...)
		}

	}
	runtime.GC()
	fmt.Printf(" GC took %s\n", utils.TimeGC())
}

func Benchmark_Experiment3Copy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 100)
		n := 0
		for j := 0; j < 20; j++ {
			n += copy(s[n:], []int{j, j + 1, j + 2, j + 3, j + 4})
		}
	}
	runtime.GC()
	fmt.Printf(" GC took %s\n", utils.TimeGC())
}
