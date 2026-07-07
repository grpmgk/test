package main

import "fmt"

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// bfsMin возвращает минимальное значение в дереве, используя обход в ширину
func bfsMin(root *Node) int {
	stack := []*Node{root}
	min := root.Value

	for len(stack) > 0 {
		curr := stack[0]
		stack = stack[1:]

		if curr.Value < min {
			min = curr.Value
		}

		if curr.Left != nil {
			stack = append(stack, curr.Left)
		}
		if curr.Right != nil {
			stack = append(stack, curr.Right)
		}
	}
	return min
}

// bfsPrint печатает все узлы по уровням (для наглядности)
func bfsPrint(root *Node) {
	if root == nil {
		return
	}
	queue := []*Node{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Printf("%d ", node.Value)
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	fmt.Println()
}

func main() {
	//        5
	//       / \
	//      3   8
	//     / \   \
	//    1   4   9
	root := &Node{5,
		&Node{3, &Node{1, nil, nil}, &Node{4, nil, nil}},
		&Node{8, nil, &Node{9, nil, nil}},
	}

	fmt.Print("BFS обход: ")
	bfsPrint(root) // 5 3 8 1 4 9

	fmt.Println("Минимальный элемент:", bfsMin(root)) // 1
}
