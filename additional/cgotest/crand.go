package cgotest

/*
#include <stdlib.h>
*/
import "C"

// Random ...
func Random() int {
	return int(C.random())
}

// Seed ...
func Seed(i int) {
	C.srandom(C.uint(i))
}