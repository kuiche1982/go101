package main

import "fmt"

type argError struct {
	a1 int
	a2 string
}

func (a *argError) Error() string {
	return fmt.Sprintf("this is an arg eror, a1: %v, a2: %v", a.a1, a.a2)
}

func raiseError() (int, error) {
	return 0, &argError{a1: 1, a2: "some error"}
}

func setDone(ch chan<- int) {
	fmt.Println("sending signal")
	ch <- 1
	fmt.Println("signal sent")
}

func waitDone(ch <-chan int) {
	fmt.Println("start wait chan signal")
	<-ch
	fmt.Println("end wait chan signal")
}
func p() {
	done := make(chan int)
	// bufferred channel as following;
	// ch := make(chan int, 10)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(argError); ok {
					fmt.Println(e)
				}
				fmt.Println("resolved the error")
			}

		}()
		setDone(done)
	}()
	waitDone(done)
}
