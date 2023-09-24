package event

import (
	"skeleton/app/event/listen"
	"skeleton/internal/variable"
)

type Event struct {
}

func (*Event) Init() error {
	err := variable.Event.Register(&listen.FooListen{})
	if err != nil {
		return err
	}
	return nil
}
