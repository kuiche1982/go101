package main

import (
	"fmt"
)

func main() {
	var yy = [...]string{"asdf"} // define a cap & len just enough array
	//yy = append(yy, "abcd")   // this line could not complie. since yy is not slice
	fmt.Printf("%t, \n%s\n", yy, yy)
	fmt.Println(len(yy), cap(yy))

}
