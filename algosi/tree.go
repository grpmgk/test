package main

import (
	"fmt"
	"strconv"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func bulbulator3000(node *Node, s string) string {
	if node == nil {
		return s
	}
	s = bulbulator3000(node.Left, s)
	s += strconv.Itoa(node.Value)
	s = bulbulator3000(node.Left, s)
	return s
}

func main() {
	root := &Node{Value: 5}
	root.Left = &Node{Value: 3}
	root.Right = &Node{Value: 8}
	root.Left.Left = &Node{Value: 1}
	root.Left.Right = &Node{Value: 4}
	root.Right.Right = &Node{Value: 9}

	fmt.Println(bulbulator3000(root, ""))
}
