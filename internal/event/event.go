package event

import (
	"errors"
	"log"
	"reflect"
	"sync"
)

type IEvent interface {
	GetData() any
}

type IListen interface {
	Listen() []IEvent
	Process(param any)
}

var smap sync.Map

type Event struct {
	wg *sync.WaitGroup
}

func New(listens ...IListen) *Event {
	evt := &Event{
		wg: new(sync.WaitGroup),
	}
	if len(listens) == 0 {
		return evt
	}
	if err := evt.Register(listens...); err != nil {
		log.Fatalf("event register: %s", err)
	}
	return evt
}

func (e *Event) Register(listens ...IListen) error {
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

func (e *Event) on(event IEvent, listen IListen) {
	name := reflect.TypeOf(event).String()
	evts, ok := smap.Load(name)
	var handlers []IListen
	if ok {
		handlers = evts.([]IListen)
	}
	if len(handlers) == 0 {
		smap.Store(name, []IListen{listen})
		return
	}
	handlers = append(handlers, listen)
	smap.Store(name, handlers)
}

func (e *Event) Dispatch(event IEvent) error {
	name := reflect.TypeOf(event).String()
	handlers, ok := smap.Load(name)
	if !ok {
		return errors.New("event not registered")
	}
	return e.exec(handlers.([]IListen), event, false)
}

func (e *Event) DispatchAsync(event IEvent) error {
	name := reflect.TypeOf(event).String()
	handlers, ok := smap.Load(name)
	if !ok {
		return errors.New("event not registered")
	}
	return e.exec(handlers.([]IListen), event, true)
}

func (e *Event) exec(handlers []IListen, event IEvent, async bool) error {
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
		go func() {
			defer e.wg.Done()
			execFunc.Process(param)
		}()
	}
	e.wg.Wait()
	return nil
}
