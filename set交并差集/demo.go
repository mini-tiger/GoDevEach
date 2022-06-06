package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"log"
)

type A struct {
	AA string
}

func main() {
	var a1 A = A{"1"}
	var a2 A = A{"1"}   // todo 值一样就去重
	var a3 *A = &A{"1"} // todo 指针不能去重
	kide1 := mapset.NewSet()
	kide1.Add(a1)
	kide1.Add(a2)
	kide1.Add(a1)
	kide1.Add(a3)
	fmt.Printf("kide1 len:%d,mem:%+v\n", kide1.Cardinality(), kide1.ToSlice())

	kide := mapset.NewSet()
	kide.Add("xiaorui.cc")
	kide.Add("blog.xiaorui.cc")
	kide.Add("vps.xiaorui.cc")
	kide.Add("linode.xiaorui.cc")

	special := []interface{}{"Biology", "Chemistry"}
	scienceClasses := mapset.NewSetFromSlice(special)

	address := mapset.NewSet()
	address.Add("beijing")
	address.Add("nanjing")
	address.Add("shanghai")

	bonusClasses := mapset.NewSet()
	bonusClasses.Add("Go Programming")
	bonusClasses.Add("Python Programming")

	//一个并集的运算
	allClasses := kide.Union(scienceClasses).Union(address).Union(bonusClasses)
	fmt.Printf("union:%+v\n", allClasses)

	//是否包含"Cookiing"
	fmt.Printf("是否包含:%v\n", scienceClasses.Contains("Cooking")) //false

	//两个集合的差集
	fmt.Printf("Difference:%v\n", allClasses.Difference(scienceClasses)) //Set{Music, Automotive, Go Programming, Python Programming, Cooking, English, Math, Welding}

	//两个集合的交集
	fmt.Printf("Insersect:%v\n", scienceClasses.Intersect(kide)) //Set{Biology}

	//有多少基数
	fmt.Printf("cardinality:%d\n", bonusClasses.Cardinality()) //2

	log.Printf("scienceClasses 是否是allColassese 子集:%v\n", scienceClasses.IsSubset(allClasses))
	log.Printf("scienceClasses 是否是allColassese 超集:%v\n", allClasses.IsSuperset(scienceClasses))

	num1 := mapset.NewSet()
	num1.Add(2)
	num1.Add(1)
	num2 := mapset.NewSet()
	num2.Add(1)
	num2.Add(2)

	fmt.Println(num1.Equal(num2))
}
