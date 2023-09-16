package main

import "fmt"

func adder(num1, num2 int) int {
	return num1 + num2
}

func main() {
    val := adder(2, 2)
    fmt.Println(val)
}
