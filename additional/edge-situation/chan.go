package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var ch chan int

	go func(ch chan int) {
		time.Sleep(time.Second)
		<-ch // never be able to dequeue any data,
		fmt.Printf("dequeued\n")
	}(ch) // at this time ch == nil
	go func() {
		ch = make(chan int, 1) // this chaned the ch value
		ch <- 1                // put succeeded since it's cached chan. this goroutine quite
		fmt.Printf("enqueued\n")
	}()
	c := time.Tick(1 * time.Second)
	for range c {
		fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine()) // print 2 goroutine always
	}
}
