package main

import (
	"fmt"
	"os"
)

func exit() {
	defer fmt.Print("already terminated")
	os.Exit(3)
}
