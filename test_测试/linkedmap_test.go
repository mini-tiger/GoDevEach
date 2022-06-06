package test

import (
	lm "github.com/mini-tiger/tjtools/LinkedMap"
	"strconv"
	"testing"
)

var lmm = lm.NewLinkedMap()

func loadData() {
	for i := 0; i < 200; i++ {
		lmm.Put(strconv.Itoa(i), i)
	}

	for i := 0; i < 100; i++ {
		lmm.Remove(strconv.Itoa(i))
	}
}
func Benchmark_Linkmap(b *testing.B) {
	//a := lm.NewLinkedMap()

	loadData()

	//fmt.Println(a.Max())
	//fmt.Printf("%+v\n", a.MData)
	//fmt.Printf("%+v\n", a.MLink)
	//for _, key := range lmm.SortLinkMap() {
	//fmt.Println(key)
	//lmm.Get(key)
	//fmt.Printf("key:%s,value:%v\n", key, v)

	//}
	for i := 0; i < b.N; i++ {
		lmm.Get(strconv.Itoa(i))
	}

}
