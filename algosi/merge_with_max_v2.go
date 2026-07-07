package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

func mergeSortedLists(l1, l2 *Node) *Node {
	dummy := &Node{}
	cur := dummy

	for l1 != nil || l2 != nil {
		var val int
		if l1 == nil {
			val = l2.Val
		} else if l2 == nil {
			val = l1.Val
		} else {
			if l1.Val <= l2.Val {
				val = l1.Val
			} else {
				val = l2.Val
			}
		}

		startL1 := l1
		startL2 := l2
		count1, count2 := 0, 0

		var lastL1, lastL2 *Node

		for l1 != nil && l1.Val == val {
			count1++
			lastL1 = l1
			l1 = l1.Next
		}
		for l2 != nil && l2.Val == val {
			count2++
			lastL2 = l2
			l2 = l2.Next
		}

		if count1 >= count2 {
			cur.Next = startL1
			cur = lastL1
		} else {
			cur.Next = startL2
			cur = lastL2
		}
	}
	return dummy.Next
}

func printList(head *Node) {
	for cur := head; cur != nil; cur = cur.Next {
		fmt.Printf("%d ", cur.Val)
	}
	fmt.Println()
}

func main() {
	// Создаём первый список: 1 -> 3 -> 5
	node1 := &Node{Val: 1}
	node2 := &Node{Val: 1}
	node3 := &Node{Val: 5}
	node1.Next = node2
	node2.Next = node3

	// Создаём второй список: 2 -> 4 -> 6
	node4 := &Node{Val: 1}
	node5 := &Node{Val: 5}
	node6 := &Node{Val: 5}
	node4.Next = node5
	node5.Next = node6

	fmt.Print("Список 1: ")
	printList(node1)
	fmt.Print("Список 2: ")
	printList(node4)

	merged := mergeSortedLists(node1, node4)

	fmt.Print("Результат слияния: ")
	printList(merged)
}
