package main

import (

	"sort"
	"fmt"
)

type a []string


func (tt a)Len() int {
	return len(tt)
}

func (tt a)Less(i,j int) bool  {
	return tt[i] < tt[j]
}

func (tt a)Swap(i,j int)  {
	tt[i], tt[j] = tt[j],tt[i]
}

func main() {
	var aa a= a{"a1","z1","a1","b1"}
	var ttttt  sort.Interface
	ttttt=&aa
	sort.Stable(ttttt)
	fmt.Println(ttttt)


	//ai:=[]int{3,2,1,4}
	var ss sort.IntSlice
	ss=[]int{3,2,1,4}
	ss.Sort()
	fmt.Println(ss)

}
