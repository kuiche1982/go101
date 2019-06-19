package main

// https://stackoverflow.com/questions/30367803/why-declare-like-var-i-t-and-var-i-t-continuously-in-golang/30367931

import (
	"context"
	"errors"
	"fmt"
)

type baseFunctionClass struct {
	funcName string
	minArgs  int
	maxArgs  int
}
type someFunc struct {
	baseFunctionClass
}

type builtinFunc interface {
	// getFunction gets a function signature by the types and the counts of given arguments.
	getFunction(ctx context.Context, args ...string) ([]string, error)
}

var (
	_ builtinFunc = &someFunc{}
)

func (sf *someFunc) getFunction(ctx context.Context, args ...string) ([]string, error) {
	return args, errors.New("dummy err")
}
func main() {
	bf := &someFunc{}
	bfc, err := bf.getFunction(context.Background(), "string", "another string")
	fmt.Println(len(bfc))
	fmt.Println(bfc, err)
}
