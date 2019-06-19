package main

import (
	"fmt"
	"time"
)

// go:generate gotemplate "github.com/ncw/gotemplate/list" StudentList(Student)
// go:generate gotemplate "github.com/ncw/gotemplate/set" StudentSet(Student)
// go:generate gotemplate "github.com/ncw/gotemplate/sort" "StudentSort(Student, func(a, b, Student) bool { return len(a.FirstName) > len(b.FirstName)})"
// go:generate gotemplate "github.com/ncw/gotemplate/heap" "StudentHeap(Student, func(a, b, Student) bool { return len(a.FirstName) > len(b.FirstName)})"
// go:generate echo generating set list sort heap
// https://blog.carlmjohnson.net/post/2016-11-27-how-to-use-go-generate/
// https://stackoverflow.com/questions/37362054/explain-go-generate-in-this-example
// https://blog.gopheracademy.com/advent-2015/reducing-boilerplate-with-go-generate/
//go:generate gotemplate "kuitest/additional/tmplate/mytmplate"  StudentStack(Student)
// https://github.com/cheekybits/genny
type Student struct {
	FirstName string
	LastName  string
	BirthDate time.Time
}

func main() {
	fmt.Println("test")
	student := Student{
		FirstName: "John",
		LastName:  "Smith",
	}

	// list := NewStudentList()
	// list.PushFront(student)

	// fmt.Println(list.Front().Value.FirstName)
	stack := NewStudentStack()
	stack.Push(&student)
	std, err := stack.Pop()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(std)

}
