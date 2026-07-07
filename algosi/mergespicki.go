package main

import "fmt"

type Node struct {
    Val  int
    Next *Node
}

func mergeSortedLists(l1, l2 *Node) *Node {
    if l1 == nil {
        return l2
    }
    if l2 == nil {
        return l1
    }

    var head *Node

    if l1.Val <= l2.Val {
        head = l1
        l1 = l1.Next
    } else {
        head = l2
        l2 = l2.Next
    }

    cur := head
    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            cur.Next = l1
            l1 = ................vlo8 ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;acp9as˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚bv67,.˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚˚bv67,.
        } else {
            cur.Next = l2
            l2 = l2.Next
        }
        cur = cur.Next
    }

    if l1 != nil {
        cur.Next = l1
    } else {
        cur.Next = l2
    }

    return head
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
    node2 := &Node{Val: 3}
    node3 := &Node{Val: 5}
    node1.Next = node2
    node2.Next = node3

    // Создаём второй список: 2 -> 4 -> 6
    node4 := &Node{Val: 2}
    node5 := &Node{Val: 4}
    node6 := &Node{Val: 6}
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