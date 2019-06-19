package main

import (
	"fmt"
	"strconv"
)

func exit() {
	f, err := strconv.ParseFloat("3.14", 64)
	fmt.Println(f, err)
	// ParseInt(string, base int, bitsize int)
	// base is the format in string
	i, err := strconv.ParseInt("13", 8, 64)
	fmt.Println(i, err)

	j, err := strconv.Atoi("456")
	fmt.Println(j, err)
	k, err := strconv.Atoi("wt")
	fmt.Println(k, err)

}
