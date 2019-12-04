package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	wg := sync.WaitGroup{}
	wg.Add(10)
	<-ticker.C
	time.Sleep(time.Duration(10) * time.Second)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	for i := 0; i < 10; i++ {
		wg.Done()
		select {
		case <-ticker.C:
			{
				// no blocking operation, since not rolling add rate to the chan
				fmt.Println("I'm here")
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			}
		case <-time.After(5 * time.Second):
			{
				fmt.Println("I'm default")
			}
		}
	}
	wg.Wait()
}
