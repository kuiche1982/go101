package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.RWMutex
var count int

func main() {
	go A()
	go func() {
		time.Sleep(2 * time.Second)
		mu.Lock()
		defer mu.Unlock()
		count++
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		fmt.Println(count)
	}
}

func A() {
	mu.RLock()
	defer mu.RUnlock()
	B()
}

func B() {
	time.Sleep(5 * time.Second)
	C()
}

func C() {
	mu.RLock()
	defer mu.RUnlock()
}
