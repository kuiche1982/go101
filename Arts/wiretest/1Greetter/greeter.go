package main

import "fmt"

type Message string
type Greeter struct {
	Message Message
}
type Event struct {
	Greeter Greeter
}

func NewMessage() Message {
	return Message("hi there!")
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	msg := NewMessage()
	gre := NewGreeter(msg)
	event := NewEvent(gre)
	event.Start()
}
