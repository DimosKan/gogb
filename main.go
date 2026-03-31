package main

import "fmt"

func Add(a int, b int) int {
	return a + b
}

func main() {
	result := Add(1, 1)
	fmt.Println(result)
}
