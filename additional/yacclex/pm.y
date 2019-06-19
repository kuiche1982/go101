//go:generate goyacc -o pm.go -p Cal pm.y
%{
package main
import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)
var Result = 0
%}
%union{
    val int
}
%token NUMBER
%%
list : 
    | list exp {
        fmt.Println($$, $1, $2, "exp called")
        $$.val = $2.val
        Result = $2.val
    }
exp : NUMBER
    | exp '+' NUMBER 
    {
        fmt.Println($$, $1, $3, "add called")
        $$.val = Add($1.val, $3.val)
    }
    | exp '-' NUMBER 
    {
        fmt.Println($$, $1, $3, "minus called")
        $$.val = Minus($1.val, $3.val)
    };
%%

func Add(a, b int) int {
    return a + b
}
func Minus(a, b int) int {
    return a - b
}

type CalLex struct {
	s string
	pos int
}


func (l *CalLex) Lex(lval *CalSymType) int {
	var c rune = ' '
	for c == ' ' {
		if l.pos == len(l.s) {
			return 0
		}
        c = rune(l.s[l.pos])
        l.pos += 1
        if c == '+' || c == '-' {
            return int(c)
        }
        if unicode.IsDigit(c) {
            lval.val = int(c - '0')
            return NUMBER
        }
	}
	return 0
}

func (l *CalLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))

	for {
		var eqn string
		var ok bool

		fmt.Printf("equation: ")
		if eqn, ok = readline(fi); ok {
            Result = 0
			CalParse(&CalLex{s: eqn})
            fmt.Println(Result)
		} else {
			break
		}
	}
}

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	return s, true
}