package main

import (
	"fmt"
)

// does go lang have real 2d array like [3,4]int
func p() {
	s := make([]string, 3)
	fmt.Println("s:", s, len(s), cap(s))
	// 3 empty string
	n := make([]int, 0, 3)
	fmt.Println("n:", n, len(n), cap(n))
	// emtpy slice, no string with length 0, cap 3

	n = append(n, 1, 2, 3, 4)
	fmt.Println("n:", n, len(n), cap(n))
	// auto increase the slice cap to double of the target size. 8 here

	c := make([]int, 2)
	copy(c, n)
	fmt.Println("c:", c, len(c), cap(c))

	m := n[:]
	fmt.Println("m:=n[:]", m, len(m), cap(m))

	m = n[1:2]
	fmt.Println("m:=n[1:2]", m, len(m), cap(m))

	m = n[1:]
	fmt.Println("m:=n[1:]", m, len(m), cap(m))

	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d:", twoD)
}
