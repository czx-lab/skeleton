package listen

import (
	"fmt"
	event2 "github.com/czx-lab/skeleton/app/event/entity"
	"github.com/czx-lab/skeleton/internal/event"
)

type FooListen struct {
}

func (*FooListen) Listen() event.EventInterface {
	return &event2.FooEvent{
		Name: "测试",
	}
}

func (*FooListen) Process(data any) (any, error) {
	return fmt.Sprintf("%v --> %s", data, "exec FooListen.Process"), nil
}
