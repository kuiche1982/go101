// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"kuitest/Arts/wiretest/lib"
)

// Injectors from wire.go:

func InitEvent(param string) (lib.Event, error) {
	message := lib.NewMessage(param)
	greeter := lib.NewGreeter(message)
	event, err := lib.NewEvent(greeter)
	if err != nil {
		return lib.Event{}, err
	}
	return event, nil
}

func InitEvent2(param string) (lib.Event2, error) {
	message := _wireMessageValue
	greeter := lib.Greeter{
		Message: message,
	}
	event2 := lib.Event2{
		Greeter: greeter,
	}
	return event2, nil
}

var (
	_wireMessageValue = lib.Message("some message")
)

// wire.go:

var wireSet = wire.NewSet(wire.Value(lib.Message("some message")), wire.Struct(new(lib.Greeter), "Message"), wire.Struct(new(lib.Event), "Greeter"), wire.Struct(new(lib.Event2), "Greeter"))