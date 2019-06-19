package main

import (
	"fmt"
)

func p() {
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i += 1
	}

	for j := 7; j < 9; j++ {
		fmt.Println(j)
	}

	for {
		fmt.Println("loop ", i)
		i++
		if i > 15 {
			break
		}
	}

	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}
}
