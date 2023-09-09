package event

import (
	"errors"
	"sync"
)

type InitInterface interface {
	Init() error
}

type EventInterface interface {
	EventName() string
	GetData() any
}

type Interface interface {
	Listen() EventInterface
	Process(param any) (any, error)
}

var smap sync.Map

type Event struct {
}

func New() *Event {
	return &Event{}
}

func (e *Event) Register(event Interface) error {
	if _, ok := smap.Load(event.Listen().EventName()); ok {
		return errors.New("event registered")
	}
	smap.Store(event.Listen().EventName(), event)
	return nil
}

func (e *Event) Dispatch(event EventInterface) (any, error) {
	eventClass, ok := smap.Load(event.EventName())
	if !ok {
		return nil, errors.New("event not registered")
	}
	return e.exec(eventClass.(Interface), event)
}

func (e *Event) exec(eventClass Interface, event EventInterface) (any, error) {
	param := event.GetData()
	return eventClass.Process(param)
}
