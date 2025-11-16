package main

import "fmt"

func main() {
	result := div(10, 0)

	fmt.Println(result)
}

func div(num int, den int) int {
	if den == 0 {
		fmt.Println("Division by zero")
		return 0
	}

	return num / den
}
