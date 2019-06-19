package main

import (
	"fmt"
)

// does go lang have real 2d array like [3,4]int
func p() {
	m := make(map[string]int)
	m["1"] = 1
	m["2"] = 2

	for key, val := range m {
		fmt.Println("range map: ", key, val)
	}

	ch := make(chan int)
	go func() {
		defer close(ch)
		// fatal error: all goroutines are asleep - deadlock
		for i := 0; i < 3; i++ {
			ch <- i
		}
		// must close, otherwise range will be blocked
	}()
	for v := range ch {
		fmt.Print(" ,", v)
	}
	fmt.Println()
	t := []int{5, 4, 3, 2, 1}
	for idx, val := range t {
		fmt.Println("index:", idx, "value:", val)
	}

}
