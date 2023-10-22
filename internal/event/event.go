package event

import (
	"errors"
	"reflect"
	"sync"
)

type InitInterface interface {
	Init() error
}

type EventInterface interface {
	GetData() any
}

type Interface interface {
	Listen() []EventInterface
	Process(param any)
}

var smap sync.Map

type Event struct {
	wg sync.WaitGroup
}

func New() *Event {
	return &Event{
		wg: sync.WaitGroup{},
	}
}

func (e *Event) Register(listens ...Interface) error {
	for _, listen := range listens {
		events := listen.Listen()
		if len(events) == 0 {
			return errors.New("listening events cannot be empty")
		}
		for _, event := range events {
			e.on(event, listen)
		}
	}
	return nil
}

func (e *Event) on(event EventInterface, listen Interface) {
	name := reflect.TypeOf(event).String()
	evts, ok := smap.Load(name)
	var handlers []Interface
	if ok {
		handlers = evts.([]Interface)
	}
	if len(handlers) == 0 {
		smap.Store(name, []Interface{listen})
		return
	}
	handlers = append(handlers, listen)
	smap.Store(name, handlers)
}

func (e *Event) Dispatch(event EventInterface) error {
	name := reflect.TypeOf(event).String()
	handlers, ok := smap.Load(name)
	if !ok {
		return errors.New("event not registered")
	}
	return e.exec(handlers.([]Interface), event, false)
}

func (e *Event) DispatchAsync(event EventInterface) error {
	name := reflect.TypeOf(event).String()
	handlers, ok := smap.Load(name)
	if !ok {
		return errors.New("event not registered")
	}
	return e.exec(handlers.([]Interface), event, true)
}

func (e *Event) exec(handlers []Interface, event EventInterface, async bool) error {
	param := event.GetData()
	if !async {
		for _, handler := range handlers {
			handler.Process(param)
		}
		return nil
	}
	for _, handler := range handlers {
		e.wg.Add(1)
		execFunc := handler
		go func(wg *sync.WaitGroup) {
			defer e.wg.Done()
			execFunc.Process(param)
		}(&e.wg)
	}
	e.wg.Wait()
	return nil
}
