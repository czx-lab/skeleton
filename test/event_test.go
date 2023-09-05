package test

import (
	"fmt"
	"testing"

	"github.com/czx-lab/skeleton/internal/event"
)

type DemoEventEntry struct {
}

func (*DemoEventEntry) EventName() string {
	return "demo-event"
}

func (*DemoEventEntry) GetData() any {
	return 12
}

type DemoEventEntryListen struct {
}

func (*DemoEventEntryListen) Listen() event.EventInterface {
	return &DemoEventEntry{}
}

func (*DemoEventEntryListen) Process(data any) (any, error) {
	fmt.Println(data)
	return "success", nil
}

func TestEvent(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
		}
	}()
	eventClass := event.New()
	if err := eventClass.Register(&DemoEventEntryListen{}); err != nil {
		t.Log(err)
	}
	if result, err := eventClass.Dispatch(&DemoEventEntry{}); err != nil {
		t.Log(err)
	} else {
		t.Log(result)
	}
}
