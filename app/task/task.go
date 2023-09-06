package task

import (
	"github.com/czx-lab/skeleton/internal/crontab"
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
