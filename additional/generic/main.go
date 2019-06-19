package main

import "fmt"

// type List(type T) []T
// func (l List(type T)) Push(v type T) List(type T) {
// 	return append(l, v)
// }
// well here are some comments here
// https://appliedgo.net/generics/
// 1. consider copy & past code
// 2. consider use interface
// 3. use type assertion out, ok := v.(sometype)
// 4. use reflection
// 5. use code generator ( to generate the array / list code)
// I was using C# before switch to golang,
// look back to year 2004, there's only a Collection, List which accept object parameter
// at that time if we want something like List<T> we will need inhirement the array / map
// then offer necessary function to convert the value back.
// it should be the same case in golang today.
// well for the template type T, you can use interface instead.
type SomeType interface {
	DoSomething() int
}

type SomeTypeImp struct {
	simvar int
}

func (sti SomeTypeImp) DoSomething() int {
	fmt.Println(sti.simvar)
	return sti.simvar
}
func callDo(ins interface {
	DoSomething() int
}) {
	ins.DoSomething()
}
func main() {
	var v SomeType
	fmt.Println(v)
	v = &SomeTypeImp{
		simvar: 10,
	}
	fmt.Println(v)
	callDo(v)
}
