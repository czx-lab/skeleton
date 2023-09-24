package task

import (
	"skeleton/internal/crontab"
)

type Task struct {
}

var _ crontab.TaskInterface = (*Task)(nil)

func New() crontab.TaskInterface {
	return &Task{}
}

func (*Task) Tasks() crontab.Tasks {
	return []crontab.Interface{&DemoTask{}}
}
