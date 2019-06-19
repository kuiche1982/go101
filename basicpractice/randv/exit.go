package main

import (
	"fmt"
	"math/rand"
	"time"
)

func outputRand(r *rand.Rand) {
	for i := 0; i < 10; i++ {
		fmt.Print("i:", r.Intn(100), "\t")
	}
	fmt.Println()

}
func exit() {
	for i := 0; i < 10; i++ {
		fmt.Print("m:", rand.Intn(100), "\t")
	}
	fmt.Println()

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	outputRand(r)
}
