package main

import (
	"fmt"
)

func worker(i int, jobs <-chan int, out chan<- int) {
	for j := range jobs {
		fmt.Println("get item from jobs: ", j)
		out <- j + 1
		fmt.Println("output j +1 to out chan, in thread:", i)
	}
}

func p() {
	ch1 := make(chan int, 100)
	ch2 := make(chan int, 100)

	for i := 0; i < 3; i++ {
		go worker(i, ch1, ch2)
	}

	for j := 0; j < 5; j++ {
		ch1 <- j
	}
	close(ch1)
	//close(ch2)
	for k := 0; k < 5; k++ {
		t := <-ch2
		fmt.Println("receive result ", k, " value:", t)
	}
}
