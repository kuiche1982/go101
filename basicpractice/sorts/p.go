package main

import (
	"fmt"
	"sort"
)

type Ppl struct {
	name string
	age  int
}
type PplArray []Ppl

func (p PplArray) Len() int {
	return len(p)
}

func (p PplArray) Less(i, j int) bool {
	return p[i].age < p[j].age
}
func (p PplArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// sort.Ints([]int)
// sort.Strings([]string)

func p() {
	ppls := PplArray{
		Ppl{
			name: "3",
			age:  3,
		},
		Ppl{
			name: "4",
			age:  4,
		}, Ppl{
			name: "1",
			age:  1,
		},
		Ppl{
			name: "2",
			age:  2,
		},
	}

	fmt.Println(ppls)
	sort.Sort(ppls)
	fmt.Println(ppls)
}
