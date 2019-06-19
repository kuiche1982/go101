package main

import (
	"fmt"
	"time"
)

func p(i int) string {
	var r string
	fmt.Println("Write ", i, " as ")
	switch i {
	case 1:
		r = "one"
	case 2:
		r = "two"
	case 3:
		fallthrough
	case 4:
		r = "unknown"
	default:
		r = "defualt"
	}
	fmt.Println(r)
	return r
}

func q() {
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("week end")
	default:
		fmt.Println("work day")
	}
}

func w() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("morning")
	default:
		fmt.Println("afternoon")
	}
}

var whatAmI = func(i interface{}) {
	switch t := i.(type) {
	case bool:
		fmt.Println("I'm bool")
	case int:
		fmt.Println("I'm an int")
	default:
		fmt.Printf("Don't know type %T\n", t)
	}
}
