%{  
    package main  
    import (  
        "fmt"  
        "os"  
        "io/ioutil"  
        "flag"  
        "bufio"  
        )  
    var fi *bufio.Reader  
    var peekrune int  
    var data []byte  
    var linep int = 0  
    var finalval float64 = 0  
%}  
  
%union  
{  
    vvar   string;  
    numval float64;  
}  
  
%token NUMBER  
%token OP  
  
%%  
expr:   
         expr1  
| expr '+' expr1  
{  
    $$.numval = $1.numval + $3.numval  
    finalval = $$.numval  
}  
| expr '-' expr1  
{  
    $$.numval = $1.numval - $3.numval  
    finalval = $$.numval  
}  
  
expr1:  
         NUMBER  
|        expr1 '*' NUMBER  
{  
    $$.numval = $1.numval * $3.numval  
    finalval = $$.numval  
}  
|       expr1 '/' NUMBER  
{  
    $$.numval = $1.numval /   
        $3.numval  
    finalval = $$.numval  
}  
  
%%  
  
func getrune() int {  
    if linep >= len(data) {  
        return 0  
    }  
    c := data[linep]  
      
    return int(c)  
}  
  
func next() {  
    linep++  
}  
  
func getnumber(c int) int {  
    var n int = 0  
    for ;c>='0' && c <= '9'; {  
        n += (c - '0')  
        next()  
        c = getrune()  
    }  
    yylval.numval = float64(n)  
    return NUMBER  
      
}  
  
func readblank() {  
    var c int  
    for c = getrune(); c == ' '; {  
        next()  
        c = getrune()  
    }  
}  
  
func Lex() int {  
    var c int  
    readblank()  
    c = getrune()  
    if c >= '0' && c <= '9' {  
        return getnumber(c)  
    }  
    switch c {  
    case '+', '-', '*', '/':  
        yylval.vvar = string(c)  
        next()  
        return c  
    }  
    return c  
}  
  
func Error(s string, v ...interface{}) {  
    fmt.Printf("ERROR:%s\n", s)  
      
}  
  
func main() {  
    if flag.NArg() == 0 {  
        fmt.Printf("Usage goexpr <expr file>\n")  
        os.Exit(1)  
    }  
    file := flag.Arg(0)  
    f, err := os.Open(file)  
    if err != nil {  
        fmt.Printf("Error opening %v: %v", file, err)  
        os.Exit(2)  
    }  
    data, err = ioutil.ReadAll(f)  
    if err != nil {  
        fmt.Printf("Error reading file %v, %v\n", file, err)  
        os.Exit(3)  
    }  
  
    Parse()   // Parse the data  
    fmt.Printf("result = %g\n", finalval)  
}  