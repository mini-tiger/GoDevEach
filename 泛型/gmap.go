package main

import (
	"fmt"
	"strings"
)

//需要使用几个 any 写在函数名后面[T1 any, T2 any]
func gMap[T1 any, T2 any](arr []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(arr))
	for i, elem := range arr {
		result[i] = f(elem)
	}
	return result
}

func main() {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	squares := gMap(nums, func(elem int) int {
		return elem * elem
	})
	fmt.Println(squares) //0 1 4 9 16 25 36 49 64 81
	strs := []string{"Hao", "Chen", "MegaEase"}
	upstrs := gMap(strs, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Println(upstrs) // HAO CHEN MEGAEASE
	dict := []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	strs = gMap(nums, func(elem int) string {
		return dict[elem]
	})
	fmt.Println(strs) // 零 壹 贰 叁 肆 伍 陆 柒 捌 玖
}
