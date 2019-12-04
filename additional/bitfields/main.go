package main

// #include <stdio.h>
// #include <stdlib.h>
// #include "bitfields.c"
import "C"
import (
	"fmt"
)

func main(){
	cv := C.newBitmap()
	fmt.Printf("%b, %d\n",cv, C.getField1(cv))
	cv = C.setField1(cv)
	fmt.Printf("%b, %d\n",cv, C.getField1(cv))
	cv = C.clearField1(cv)
	fmt.Printf("%b, %d\n",cv, C.getField1(cv))
	cv = C.setField2(cv)
	fmt.Printf("%b, %d\n",cv, C.getField2(cv))
	cv = C.clearField2(cv)
	fmt.Printf("%b, %d\n",cv, C.getField2(cv))
	C.printsize()
}
