package event

import (
	"github.com/czx-lab/skeleton/app/event/listen"
	"github.com/czx-lab/skeleton/internal/variable"
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
