package main

// if your program uses any //export directives, then the C code in the comment may only include declarations (extern int f();), not definitions (int f() { return 1; }). You can use //export directives to make Go functions accessible to C code.

// #include <stdio.h>
// #include <stdlib.h>
// #include "print.c"
import "C"
import (
	"unsafe"
)

// Print calls c print function
func Print(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.fputs(cs, (*C.FILE)(C.stdout))
	C.myprint(cs)
}


func main() {
	Print("abcd")
	Print("efgh")
}

// ldd for linux 
// otool -L 