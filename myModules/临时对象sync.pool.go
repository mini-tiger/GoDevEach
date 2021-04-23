package main

import (
	"fmt"
	"sync"
	"time"
)

var bytePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 1024)
		return &b
	},
}

func main() {
	a := time.Now().Unix()
	fmt.Println(time.Now())
	// 不使用对象池, fixme 速度快，占用内存磊
	for i := 0; i < 500000000; i++{
		obj := make([]byte,1024)
		_ = obj
	}

	fmt.Println(time.Now())
	b := time.Now().Unix()
	// fixme 使用对象池,  速度慢，占用内存小,放入到 sync.pool后,只能取出一次，减小GC压力
	// xxx 适合 需要多个临时对象，使用一次 便丢弃的
	for i := 0; i < 500000000; i++{
		obj := bytePool.Get().(*[]byte) // xxx  提取一个 没有放入的值，调用new 方法，返回一个新的
		//fmt.Println(len(*obj))
		bytePool.Put(obj) //xxx 放入一个
	}
	fmt.Println(time.Now())
	c := time.Now().Unix()

	fmt.Println("without pool ", b - a, "s")
	fmt.Println("with    pool ", c - b, "s")
	Test1() // xxx
}

// without pool  34 s
// with    pool  24 s
func Test1() {
	// 初始化一个pool
	pool := &sync.Pool{
		// 默认的返回值设置，不写这个参数，默认是nil
		New: func() interface{} {
			return 0
		},
	}

	// 看一下初始的值，这里是返回0，如果不设置New函数，默认返回nil
	init := pool.Get()
	fmt.Println(init)

	// 设置一个参数1
	pool.Put(1)

	// 获取查看结果
	num := pool.Get()
	fmt.Println(num)

	// 再次获取，会发现，已经是空的了，只能返回默认的值。
	num = pool.Get()
	fmt.Println(num)
}