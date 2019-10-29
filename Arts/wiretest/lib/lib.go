package lib

import (
	"errors"
	"fmt"
	"time"
)

type Message string
type Greeter struct {
	Message Message
	Grumpy  bool
}
type Event struct {
	Greeter Greeter
}

func NewMessage(phrase string) Message {
	return Message(phrase)
}

func NewGreeter(m Message) Greeter {
	var grumpy bool
	if time.Now().Unix()%2 == 0 {
		grumpy = true
	}
	return Greeter{Message: m, Grumpy: grumpy}
}

func (g Greeter) Greet() Message {
	if g.Grumpy {
		return Message("Go away!")
	}
	return g.Message
}

func NewEvent(g Greeter) (Event, error) {
	if g.Grumpy {
		return Event{}, errors.New("could not create event: event greeter is grumpy")
	}
	return Event{Greeter: g}, nil
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

type Event2 struct {
	Greeter Greeter
}

func NewEvent2(g Greeter) (Event2, error) {
	if g.Grumpy {
		return Event2{}, errors.New("could not create event2: event greeter is grumpy")
	}
	return Event2{Greeter: g}, nil
}
