package test

import (
	"testing"
)

func TestSum(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expected := 15
	actual := Sum(numbers)
	actual1 := Sum1(numbers)
	if actual != expected {
		t.Errorf("Expected the sum of %v to be %d but instead got %d!", numbers, expected, actual)
	}
	if actual1 != expected {
		t.Errorf("Expected the sum1 of %v to be %d but instead got %d!", numbers, expected, actual)
	}
}

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
