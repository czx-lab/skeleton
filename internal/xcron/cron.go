package xcron

import "github.com/robfig/cron/v3"

type ICron interface {
	Rule() string
	Execute() func()
}

type Cron struct {
	instance *cron.Cron
}

func New() *Cron {
	return &Cron{
		instance: cron.New(),
	}
}

func (c *Cron) AddFunc(cmd ...ICron) {
	if len(cmd) == 0 {
		return
	}
	for _, job := range cmd {
		c.instance.AddJob(job.Rule(), cron.FuncJob(job.Execute()))
	}
}

func (c *Cron) Start() {
	c.instance.Start()
}
