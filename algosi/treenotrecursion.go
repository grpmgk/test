package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func inOrderIterative(root *Node) string {
	var sb strings.Builder
	stack := []*Node{}
	curr := root

	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		sb.WriteString(strconv.Itoa(curr.Value))

		curr = curr.Right
	}
	return sb.String()
}

func main() {
	// Дерево:
	//      a(1)
	//     /    \
	//   b(2)   e(5)
	//  /  \    /  \
	// c(3) d(4) f(6) g(7)

	root := &Node{1,
		&Node{2,
			&Node{3, nil, nil},
			&Node{4, nil, nil},
		},
		&Node{5,
			&Node{6, nil, nil},
			&Node{7, nil, nil},
		},
	}

	fmt.Println(inOrderIterative(root)) // 3 2 4 1 6 5 7
}
