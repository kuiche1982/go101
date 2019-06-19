package main

import (
	"fmt"
	"sync"
)

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	var ch1, ch2 chan int
	ch1 = make(chan int)
	ch2 = make(chan int)
	var loopfun = func(ch chan<- int, startvalue int) {
		defer close(ch)
		for i := 0; i < 5; i++ {
			iv := i + startvalue
			ch <- iv
		}
	}

	go loopfun(ch1, 0)
	go loopfun(ch2, 10)
	for a := range merge(ch1, ch2) {
		fmt.Print(";", a)
	}
	fmt.Println()
}
