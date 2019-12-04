//+build wireinject

package main

import (
	"github.com/google/wire"
)

func InitEvent() Event {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}
}
