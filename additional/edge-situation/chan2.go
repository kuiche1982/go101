package main

import "fmt"

func main() {
	var ch chan int // = make(chan int)
	var count int
	go func() {
		ch <- 1
	}()
	go func() {
		count++
		close(ch) // panic, close of nil chan
	}()
	<-ch
	fmt.Println(count)
}
