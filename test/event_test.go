package test

import (
	"fmt"
	"testing"

	"skeleton/internal/event"
)

type DemoEvent1 struct {
}

func (*DemoEvent1) GetData() any {
	return "DemoEvent1"
}

type DemoEvent2 struct {
}

func (*DemoEvent2) GetData() any {
	return "DemoEvent2"
}

type DemoEventEntryListen1 struct {
}

func (*DemoEventEntryListen1) Listen() []event.EventInterface {
	return []event.EventInterface{
		&DemoEvent2{},
		&DemoEvent1{},
	}
}

func (*DemoEventEntryListen1) Process(data any) {
	fmt.Printf("Execute listener => %s function Process, event => %v\n", "DemoEventEntryListen1", data)
}

type DemoEventEntryListen2 struct {
}

func (*DemoEventEntryListen2) Listen() []event.EventInterface {
	return []event.EventInterface{
		&DemoEvent1{},
	}
}

func (*DemoEventEntryListen2) Process(data any) {
	fmt.Printf("Execute listener => %s function Process, event => %v\n", "DemoEventEntryListen2", data)
}

func TestEvent(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
		}
	}()
	evt := event.New()
	if err := evt.Register(&DemoEventEntryListen1{}, &DemoEventEntryListen2{}); err != nil {
		t.Log(err)
	}

	// sync
	if err := evt.Dispatch(&DemoEvent2{}); err != nil {
		t.Log(err)
	}

	// async
	if err := evt.DispatchAsync(&DemoEvent1{}); err != nil {
		t.Log(err)
	}
}
