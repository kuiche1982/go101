package main

import (
	"context"
	"fmt"
	"sync"

	"gitlab.myteksi.net/gophers/go/commons/util/parallel/gconcurrent"
)

func main() {
	usegcon()
	// var v []int
	// v = append(v, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	// ch := make(chan string, len(v))
	// wg := sync.WaitGroup{}
	// wg.Add(len(v))
	// for _, id := range v {
	// 	someid := id
	// 	go func() {
	// 		defer func() {
	// 			if r := recover(); r != nil {
	// 				e, ok := r.(error)
	// 				if ok {
	// 					fmt.Println(e)
	// 				}
	// 			}
	// 			wg.Done()
	// 		}()
	// 		if someid == 3 {
	// 			panic("Make some error")
	// 		}
	// 		rvalue := fmt.Sprintf("%v", someid*2)
	// 		ch <- rvalue
	// 	}()
	// }

	// wg.Wait()
	// close(ch)
	// for val := range ch {
	// 	fmt.Print(val, "\t")
	// }
	// fmt.Println()
}

func usegcon() {
	ctx := context.Background()
	var v []int
	v = append(v, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	ch := make(chan string, len(v))
	// if we declare ch as make(chan string), then
	// every send and receive will be block ops
	// the wg.Wait will never end.
	//defer close(ch)
	wg := sync.WaitGroup{}
	wg.Add(len(v))
	for _, id := range v {
		someid := id
		gconcurrent.Go(ctx, "abcd", func(i context.Context) error {
			defer wg.Done()
			if someid == 3 {
				panic("Make some error")
			}
			rvalue := fmt.Sprintf("%v", someid*2)
			ch <- rvalue
			return nil
		})
	}
	wg.Wait()
	close(ch)

	for val := range ch {
		fmt.Println(val, "\t")
	}
}
