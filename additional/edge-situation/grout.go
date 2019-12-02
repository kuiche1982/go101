package main

import (
	"fmt"
	"sync"
	"time"
)

type T struct {
	V int
}

func (t *T) Incr(wg *sync.WaitGroup) {
	t.V++
	wg.Done()
}
func (t *T) Print() {
	time.Sleep(1e9)
	fmt.Print(t.V)
}
func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	var ts = make([]T, 10)
	for i := 0; i < 10; i++ {
		ts[i] = T{i}
	}
	// for i := 0; i < 10; i++ {
	// 	fmt.Print(ts[i].V)
	// }
	// for _, t := range ts {
	// 	s := t
	// 	go s.Incr(&wg)
	// }
	// wg.Wait()
	// for i := 0; i < 10; i++ {
	// 	fmt.Print(ts[i].V)
	// }
	for _, t := range ts {
		// s := t
		go t.Print()
	}
	time.Sleep(5 * time.Second)
}

// A: 输出12345678910
// B: 输出0123456789
// C: 输出9999999999
// D: 输出10101010101010101010
