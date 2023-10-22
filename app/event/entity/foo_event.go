package entity

import "fmt"

type FooEvent struct {
	Name string
}

func (f *FooEvent) GetData() any {
	return fmt.Sprintf("FooEvent.Name = %s --> %s", f.Name, "exec FooEvent.GetData")
}
