package main

import "fmt"

func main() {
	m := make(map[string]string)
	mv, ok := m["notexists"]
	fmt.Println(mv, ok, len(mv))
}
