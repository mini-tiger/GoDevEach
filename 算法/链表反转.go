package main

import "fmt"

// 定义一个简单的链表节点结构
type ListNode struct {
	Val  int
	Next *ListNode
}

// 反转链表函数
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	current := head
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}
	return prev
}

// 创建链表函数
func createList(nums []int) *ListNode {
	var head *ListNode
	var tail *ListNode
	for _, num := range nums {
		node := &ListNode{Val: num}
		if head == nil {
			head = node
			tail = node
		} else {
			tail.Next = node
			tail = node
		}
	}
	return head
}

// 打印链表函数
func printList(head *ListNode) {
	current := head
	for current != nil {
		fmt.Printf("%d ", current.Val)
		current = current.Next
	}
	fmt.Println()
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	head := createList(nums)

	fmt.Println("原始链表：")
	printList(head)

	// 反转链表
	reversedHead := reverseList(head)

	fmt.Println("反转后的链表：")
	printList(reversedHead)
}
