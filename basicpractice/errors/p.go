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
	_, err := raiseError()
	if e, ok := err.(*argError); ok {
		fmt.Println(" could catch and handle the error", e)
	} else {
		fmt.Println("could not handl the error, should return or panic", e)
	}
}
