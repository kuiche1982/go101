//go:generate go build -buildmode=plugin plg.go
package main

import "fmt"

var V int

func F() { fmt.Printf("Hello, number %d\n", V) }
