package utils

import (
	"runtime"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2022/7/3
 * @Desc: slice.go
**/
func TimeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}
func GenSlice() []int {
	s := make([]int, 1023)
	for i := 0; i < 1023; i++ {
		s[i] = i
	}

	return s
}

func GenSliceAppend() []int {
	s := make([]int, 0)
	//for i := 0; i < 10; i++ {
	//	s = append(s, i)
	//}

	return s
}
