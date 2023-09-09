package entity

import "fmt"

type FooEvent struct {
	Name string
}

func (*FooEvent) EventName() string {
	return "foo-event"
}

func (f *FooEvent) GetData() any {
	return fmt.Sprintf("FooEvent.Name = %s --> %s", f.Name, "exec FooEvent.GetData")
}
