package main

import (
	"fmt"
)

// does go lang have real 2d array like [3,4]int
func p() {
	m := make(map[string]int)
	m["1"] = 1
	m["2"] = 2

	fmt.Println("make map: ", m)

	n := map[string]int{}
	n["1"] = 1
	n["2"] = 2

	fmt.Println("auto init map: ", n)

	q := map[int]int{1: 1, 2: 2}
	fmt.Println("init a map: ", q)

	y, k := m["3"]
	fmt.Printf("try to get not exists key, value: %v, exists: %v\n", y, k)

	y, k = m["2"]
	fmt.Printf("try to get exists key, value: %v, exists: %v\n", y, k)

	delete(m, "2")
	y, k = m["2"]
	fmt.Printf("try to get after delete exists key, value: %v, exists: %v\n", y, k)
}
