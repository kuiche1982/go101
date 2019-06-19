package main

import (
	"fmt"
)

func p() {
	var a = "initial"
	fmt.Println(a)
	// initial \r\n

	var b, c int = 1, 2
	fmt.Println(b, c)
	// 1 2 \r\n

	var d = true
	fmt.Println(d)
	// true \r\n

	var e int
	fmt.Println(e)
	// 0 \r\n
}
