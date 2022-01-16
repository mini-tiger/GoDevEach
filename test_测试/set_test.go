package test

import (
	"testing"
)

func Test_set(t *testing.T) {
	a := set1()
	b := set2()
	if len(a) != len(b) {
		t.Fail()
	}
}

func set1() []int {
	var x []int = []int{1, 2, 23, 4, 5, 1, 23, 4, 5, 6}
	set := make(map[int]struct{}, 0)
	setsl := make([]int, 0)
	for _, v := range x {
		set[v] = struct{}{}
	}

	for k, _ := range set {
		setsl = append(setsl, k)
	}
	return setsl
}

func set2() []int {
	var x []int = []int{1, 2, 23, 4, 5, 1, 23, 4, 5, 6}
	set := make(map[int]struct{}, 0)
	setsl := make([]int, 0)
	for _, v := range x {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		setsl = append(setsl, v)
	}

	return setsl
}

func Benchmark_Set1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = set1()
	}
}

func Benchmark_Set2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = set2()
	}
}
