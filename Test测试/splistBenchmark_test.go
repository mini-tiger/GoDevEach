package main

import (
	"Testdemo/tools"
	"testing"
)

func Benchmark_Split(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tools.Split("abcd", "bc")
	}
}

// xxx 压力测试  go test -v -bench=. -run=more splistBenchmark_test.go
