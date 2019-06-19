package main

import "fmt"

func plus(a, b int) int {
	return a + b
}

var pplus = func(a, b, c int) int {
	return a + b + c
}

func p() {
	res := plus(1, 2)
	fmt.Println("1+2=", res)
	res = pplus(1, 2, 3)
	fmt.Println("1+2+3=", res)
}

func mrvalue(a int) (int, int) {
	return a, a + a
}
