package main

import "fmt"

func main()  {
	fmt.Println(twoSum([]int{2,7,9,10,11},9)) // 返回下标
}
func twoSum(nums []int, target int) []int {
	var rint []int = make([]int,2)

	for index,num := range nums{
		for cindex,cnum := range nums[index+1:]{
			if num + cnum == target {
				rint[0]=index
				rint[1]=index+1+cindex
			}
		}
	}
	return rint
}