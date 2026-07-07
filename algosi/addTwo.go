package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

func reverse(head *Node) *Node {
	var prev *Node
	for head != nil {
		next := head.Next
		head.Next = prev
		prev = head
		head = next
	}
	return prev
}

func Summing(l1, l2 *Node) *Node {
	res := &Node{}
	cur := res
	carry := 0
	sum := 0
	for l1 != nil || l2 != nil || carry != 0 {
		sum = carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		carry = sum / 10
		cur.Next = &Node{Val: sum % 10}
		cur = cur.Next
	}
	return res.Next
}
func addTwoNumbers(l1, l2 *Node) *Node {
	Num1 := reverse(l1)
	Num2 := reverse(l2)

	res := Summing(Num1, Num2)
	return reverse(res)
}
func main() {
	// Число 12931: 1 -> 2 -> 9 -> 3 -> 1
	l1 := &Node{1, &Node{2, &Node{9, &Node{3, &Node{1, nil}}}}}

	// Число 2745: 2 -> 7 -> 4 -> 5
	l2 := &Node{2, &Node{7, &Node{4, &Node{5, nil}}}}

	result := addTwoNumbers(l1, l2)

	// Печать результата
	for result != nil {
		fmt.Print(result.Val)
		if result.Next != nil {
			fmt.Print(" -> ")
		}
		result = result.Next
	}
	fmt.Println()
}
