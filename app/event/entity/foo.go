package entity

import "fmt"

type Foo struct {
	Name string
}

func (f *Foo) GetData() any {
	return fmt.Sprintf("FooEvent.Name = %s --> %s", f.Name, "exec FooEvent.GetData")
}
