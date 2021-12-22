package test

import (
	"testing"
	"unsafe"
)

type user struct {
	name string
	age  int
}

func uuA1() *user {
	u := new(user)
	//fmt.Println(*u)

	pName := (*string)(unsafe.Pointer(u))
	*pName = "张三"

	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Offsetof(u.age)))
	*pAge = 20
	return u
}

func uuB1() *user {
	u := new(user)
	u.name = "张三"
	u.age = 20
	return u
}

func Benchmark_uu1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuA1()
	}
}

func Benchmark_uu2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuB1()
	}
}
