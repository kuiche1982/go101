package main

import (
	"fmt"
)

func p() {
	var a [5]int
	fmt.Println("emp:", a)

	a[4] = 100
	fmt.Println("set:", a)
	fmt.Println("get:", a[4])

	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("dcl : ", b)

	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}

	fmt.Println("2d: ", twoD)
	fmt.Println("len 2d:", len(twoD))
	fmt.Println("len 2d[0]:", len(twoD[0]))
	fmt.Println("cap 2d:", cap(twoD))
	fmt.Println("cap 2d[0]:", cap(twoD[0]))
}
