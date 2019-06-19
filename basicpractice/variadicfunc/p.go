package main

func plus(sums ...int) int {
	var rvalue int
	for _, v := range sums {
		rvalue += v
	}
	return rvalue
}
