package listen

import (
	"fmt"
	"czx/app/event/entity"
	"czx/internal/event"
)

type Foo struct {
}

var _ event.IListen = (*Foo)(nil)

func (*Foo) Listen() []event.IEvent {
	return []event.IEvent{
		&entity.Foo{},
	}
}

func (*Foo) Process(data any) {
	fmt.Printf("%v --> %s \n", data, "exec FooListen.Process")
}
