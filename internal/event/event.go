package event

type Interface interface {
	Listen()
	Process()
}

type Event struct {
}
