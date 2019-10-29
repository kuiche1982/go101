package main

import (
	"fmt"
	"os"
)

func main() {
	e, err := InitEvent("some init string")
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}
	e.Start()

	e2, err := InitEvent("some init string")
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}
	e2.Start()
}
