package main 

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

//export MyFunc
func MyFunc(a int, b string) int64 {
	fmt.Println(a, b)
	return int64(a)
}

//export MyFunction2
func MyFunction2() (int64, *C.char) {
	cs := C.CString("abcd")
	defer C.free(unsafe.Pointer(cs))
	return 1, cs
}


func main() {
	
}
//// go build -buildmode=c-shared -o cgotest.so cexport.go
//// go build -buildmode=c-archive -o cgotest.a cexport.go