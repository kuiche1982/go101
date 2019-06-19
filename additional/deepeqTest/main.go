package main

import (
	"fmt"

	"reflect"
)

type aimp struct {
	i int
}

type bimp struct {
	i int
	j int
}

func main() {
	ai := &aimp{
		i: 1,
	}
	bi := &aimp{
		i: 1,
	}
	fmt.Println(reflect.DeepEqual(ai, bi))
}
