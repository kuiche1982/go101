package main

import (
	"fmt"
	"strconv"
)

func test(a int, b string) bool {
	return strconv.Itoa(a) == b
}
func tests(a, b string) string {
	return a + b
}

// possible way to overload functions
func caller(ia ...interface{}) interface{} {
	pa, ok := ia[0].(int)
	pb, okb := ia[1].(string)
	pas, okas := ia[0].(string)
	if ok && okb {
		return test(pa, pb)
	}

	if okas && okb {
		return tests(pas, pb)
	}
	panic("could not cast to target parameter list")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error))
		}
	}()
	fmt.Println(caller(1, "1"))
	fmt.Println(caller("2", "2"))
}
