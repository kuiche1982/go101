package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func p() {
	var state = make(map[int]int)
	var mutex = &sync.Mutex{}

	var readOps uint64
	var writeOps uint64
	//var ops2 uint64
	for i := 0; i < 100; i++ {
		go func() {

			//ops2++
			total := 0
			for {
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)

				time.Sleep(time.Millisecond)
			}

		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)

				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)

				time.Sleep(time.Millisecond)
			}

		}()
	}
	time.Sleep(time.Second)
	opsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("read ops:", opsFinal)

	opsFinal = atomic.LoadUint64(&writeOps)
	fmt.Println("write ops:", opsFinal)

	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
	//fmt.Println("ops:", ops2)
	// go test --race
	// WARNING: DATA RACE
}
