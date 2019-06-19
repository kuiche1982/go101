package main

import "fmt"

// does go lang have real 2d array like [3,4]int
func closurefunc(nums ...int) {
	addd := getFun()
	for _, v := range nums {
		rv := addd(v)
		fmt.Println("adding closure:", v, " sum:", rv)
	}
}

func getFun() func(int) int {
	var sum int
	return func(i int) int {
		sum += i
		return sum
	}
}
