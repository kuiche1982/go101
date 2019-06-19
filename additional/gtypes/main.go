package main

import (
	"fmt"
	"runtime"

	"reflect"
)

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	fmt.Println(runtime.Caller(0))
	return t, k
}

func main() {
	a := 1
	pa := &a
	ma := map[string]bool{}
	fmt.Println(typeAndKind(a))
	fmt.Println(typeAndKind(pa))
	fmt.Println(typeAndKind(ma))
}
