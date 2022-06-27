package test

import (
	"fmt"
	"testing"
)

/**
 * @Author: Tao Jun
 * @Since: 2022/6/27
 * @Desc: arr_test.go
**/

func Benchmark_ArrLoop1(b *testing.B) {
	var arr1 []int = GenerateArr()
	for i := 0; i < b.N; i++ {
		ArrLoop1(arr1)
	}
}

func Benchmark_ArrLoop2(b *testing.B) {
	var arr1 []int = GenerateArr()
	for i := 0; i < b.N; i++ {
		ArrLoop2(arr1)
	}
}

func GenerateArr() []int {
	//var arr1 []int = make([]int, 1000)
	var arr1 []int
	for i := 0; i < 1000; i++ {
		arr1 = append(arr1, i)
		//arr1[i] = i
	}
	return arr1
}

func ArrLoop1(arr1 []int) {
	for i := 0; i < len(arr1); i++ {
		fmt.Sprint(arr1[i])
	}
}

func ArrLoop2(arr1 []int) {
	for _, i := range arr1 {
		fmt.Sprint(i)
	}
}
