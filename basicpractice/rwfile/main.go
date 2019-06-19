package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Print("already terminated")
	os.Exit(3)
}
