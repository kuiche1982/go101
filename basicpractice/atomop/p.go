package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func p() {
	var ops uint64
	//var ops2 uint64
	for i := 0; i < 50; i++ {
		go func() {
			//ops2++
			atomic.AddUint64(&ops, 1)

			time.Sleep(time.Millisecond)
		}()
	}

	time.Sleep(time.Second)
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
	//fmt.Println("ops:", ops2)
	// go test --race
	// WARNING: DATA RACE
}
