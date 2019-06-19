package main

import "fmt"

type areaer interface {
	area() float64
}

type perimer interface {
	perim() float64
}
type geometry interface {
	areaer
	perimer
}

type rect struct {
	w, h float64
}

type circle struct {
	r float64
}

func (r rect) area() float64 {
	return r.w * r.h
}

func (r rect) perim() float64 {
	return 2*r.w + 2*r.h
}

// func (r circle) area() float64 {
// 	return 3.14 * r.r * r.r
// }

func (r circle) perim() float64 {
	return 3.14 * r.r * 2
}

func measure(g interface{}) {
	fmt.Println(g)
	if ag, ok := g.(areaer); ok {
		fmt.Println(ag.area())
	}

	if pg, ok := g.(perimer); ok {
		fmt.Println(pg.perim())
	}
}

func p() {
	r := rect{w: 3, h: 4}
	c := circle{r: 6}
	measure(r)
	measure(c)
}
