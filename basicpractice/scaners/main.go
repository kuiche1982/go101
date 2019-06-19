package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	fmt.Println(args)

	wordPtr := flag.String("word", "foo", "a string")
	numPtr := flag.Int("numb", 42, "an int")
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())

	// go run main.gp  -word=opt -numb=7  -svar=flag

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		ucl := strings.ToUpper(scanner.Text())
		fmt.Println(ucl)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

}
