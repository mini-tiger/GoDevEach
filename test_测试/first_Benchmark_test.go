package test

import (
	"sync"
	"testing"
)

func Benchmark_Sum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	}
}

func Benchmark_Sum1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum1([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	}
}

func Sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func Sum1(numbers []int) int {

	sum := 0
	n := len(numbers) / 4
	var s sync.WaitGroup
	for i := 0; i <= 4; i++ {
		s.Add(1)
		b := i * n
		e := i*n + n
		go func(b, e int) {
			for _, n := range numbers[b:e] {
				sum += n
			}
			s.Done()
		}(b, e)
	}
	s.Wait()
	return sum
}
