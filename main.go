package main

import "fmt"

// 1 2
func lestnica(n int) int {

	if n <= 2 {
		return n
	}

	a, b := 1, 2
	for i := 3; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	fmt.Println(lestnica(1))
	fmt.Println(lestnica(10))
	fmt.Println(lestnica(12))
	fmt.Println(lestnica(99))

}
