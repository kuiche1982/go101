//+build wireinject

package main

import (
	lib "kuitest/Arts/wiretest/lib"

	"github.com/google/wire"
)

var wireSet = wire.NewSet(
	wire.Value(lib.Message("some message")),
	wire.Struct(new(lib.Greeter), "Message"),
	wire.Struct(new(lib.Event), "Greeter"),
	wire.Struct(new(lib.Event2), "Greeter"),
)

func InitEvent(param string) (lib.Event, error) {
	wire.Build(lib.NewEvent, lib.NewMessage, lib.NewGreeter)
	return lib.Event{}, nil
}

func InitEvent2(param string) (lib.Event2, error) {
	wire.Build(wireSet)
	return lib.Event2{}, nil
}
