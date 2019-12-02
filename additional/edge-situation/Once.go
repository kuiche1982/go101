package main

import (
	"fmt"
	"sync"
)

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if o.done == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		o.done = 1
		f()
	}
}

func main() {
	o := &Once{}
	go o.Do(func() {
		fmt.Println("doing")
	})
	go o.Do(func() {
		fmt.Println("doing")
	})
	o.Do(func() {
		fmt.Println("doing")
	})
}
