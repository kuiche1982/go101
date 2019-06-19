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
			done <- 1
		}()
		_, err := raiseError()
		panic(err)
	}()
	<-done
}
