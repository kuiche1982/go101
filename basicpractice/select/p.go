package main

import (
	"fmt"
	"time"
)

func p() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	defer close(ch1)
	defer close(ch2)
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 1
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("ch1 is awake", msg1)
		case msg2 := <-ch2:
			fmt.Println("ch2 is awake", msg2)
		}
	}
}

func timeout() {
	ch1 := make(chan int)
	defer close(ch1)

	select {
	case msg1 := <-ch1:
		fmt.Println("ch1 is awake", msg1)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout ")
	}
}

func nonblocking() {
	// should try out this, it's behavior is different
	// since there are buffer, msg is sent immediately
	// messages := make(chan string, 10)
	messages := make(chan string)
	signals := make(chan bool)

	// Here's a non-blocking receive. If a value is
	// available on `messages` then `select` will take
	// the `<-messages` `case` with that value. If not
	// it will immediately take the `default` case.
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	// A non-blocking send works similarly. Here `msg`
	// cannot be sent to the `messages` channel, because
	// the channel has no buffer and there is no receiver.
	// Therefore the `default` case is selected.
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	// We can use multiple `case`s above the `default`
	// clause to implement a multi-way non-blocking
	// select. Here we attempt non-blocking receives
	// on both `messages` and `signals`.
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func timerfun() {
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C
	fmt.Println("timer 1 expired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("timer 2 expired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("timer 2 stopped")
	}
}

func tickerfun() {
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at : ", t)
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	fmt.Println("Ticker is stopped")
}
