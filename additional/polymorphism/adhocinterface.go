package main

import (
	ic "kuitest/additional/polymorphism/income"
)

type AdHocInterface interface {
	Show() int
}

// Update 2017:
// With type assertions in Go 1.9, you can simply add = where you define the type.

// type somethingFuncy = func(int) bool
// This tells the compiler that somethingFuncy is an alternate name for func(int) bool.

// shareimprove this answer
// edited Nov 26 '17 at 0:00
// answered Nov 20 '17 at 19:10

// agarfield
// 13123
// 2
// Note that if you do this, you won't be able to add functions to the new type, it will only be an alias. You'll get an annoying cannot define new methods on non-local type error. Or at least I did when I tried this with a custom type based on string. – Ben Baron May 15 '18 at 22:27
// without following type definition
// go will not allow non local method declaration
type AAdv ic.Advertisement
type AFix ic.FixedBilling
type ATim ic.TimeAndMaterial

func (a AAdv) Show() int {
	return a.CPC * a.NoOfClicks
}

func (fb AFix) Show() int {
	return fb.BiddedAmount
}

func (tm ATim) Show() int {
	return tm.NoOfHours * tm.HourlyRate
}
