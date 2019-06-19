package main

import "fmt"

// type Int struct {
//   value int
// }

// // Imagine this is an actually complex and useful function instead.
// func (i Int) Add(j Int) Int {
//   return Int{i.value + j.value}
// }

// func main() {
//   i := Int{5}
//   j := Int{6}
//   fmt.Println(i.Add(j))
// }

// more beautiful implementation as following:

type Int int

func (i Int) Add(j Int) Int {
	return i + j
}

func main() {
	i := Int(5)
	j := Int(6)
	fmt.Println(i.Add(j))
	fmt.Println(i.Add(j) + 12)
}

// Why it matters
// The real beauty of this is actually not the ability to attach the method, although this is definitely satisfying from a conceptual point of view, bringing consistency with an approach opposite to that or Ruby, it is rather the code in the Add() method itself, and the last line in main() are actually the beauty of the thing : since our local Int type is actually just a name for the intrinsic int, it still support the builtin operators like naked ints: no explicit casts needed to use the aliased variable in an expression.
