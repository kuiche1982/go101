package main

import (
	"fmt"
	"sync"
	"time"
)

var seed int
var seedLock sync.Mutex

func newObj() interface{} {
	seedLock.Lock()
	defer seedLock.Unlock()
	time.Sleep(1 * time.Second)
	seed++
	return seed
}

func main() {
	pool := sync.Pool{
		New: newObj,
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			v := pool.Get()
			fmt.Println(seq, seed, v)
			pool.Put(v)
		}(i)
		if i%5 == 0 {
			time.Sleep(10 * time.Second)
		}
	}
	wg.Wait()
}
