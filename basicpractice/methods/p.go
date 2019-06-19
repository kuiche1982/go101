package main

import "fmt"

//Walker walker interface
type Walker interface {
	walk()
}

// Animal animal class
type Animal struct {
	kind string
}

// Ppl people class
type Ppl struct {
	Animal
	name string
}

// Walk .....
func (a *Animal) Walk() {
	fmt.Println("animal is walking:", a.kind)
}

// Walk .....
func (a *Ppl) Walk() {
	fmt.Println("ppl is walking: ", a.kind, " ", a.name)
}

func p() {
	ppl := &Ppl{}
	ppl.kind = "ppl"
	ppl.name = "ken"
	ppl2 := &Ppl{
		Animal: Animal{
			kind: "ppl",
		},
		name: "max",
	}
	anm := &Animal{kind: "cat"}
	anm.Walk()
	ppl.Walk()
	ppl2.Walk()

	//struct could not be cast ?
	// anm2 := Animal(*ppl2)
	// anm2.Walk()
	// I know unsafe pointer works, but ....
}
