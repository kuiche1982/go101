package main

import (
	"fmt"
)

var a = [...]string{
	3:"ABCD",
	4:"EFGH",
	// index:value, this is not map
}

/*
refer to go/token/token.go 
// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	literal_end
}
var token = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",
}
*/
func main() {
	fmt.Printf("%+v",a)
	fmt.Println()
	for i := 0; i < len(a); i++ {
		fmt.Println(i, ":", a[i])
	}
}
