package main

import (
	"fmt"
	"sync"
)

type MyMutex struct {
	count int
	sync.Mutex
}

func main() {
	var mu MyMutex
	mu.Lock()
	var mu2 = mu // this clause copied the lock status of mu.Mutex ( locked = 1)
	mu.count++
	mu.Unlock()
	mu2.Lock() // wait in the queue and never get unlocked
	mu2.count++
	mu2.Unlock()
	fmt.Println(mu.count, mu2.count)
}
