package cron

import (
	"fmt"
	"czx/internal/xcron"
)

type FooTask struct {
}

var _ xcron.ICron = (*FooTask)(nil)

func (d *FooTask) Rule() string {
	return "* * * * *"
}

func (d *FooTask) Execute() func() {
	return func() {
		fmt.Println("demo-task exec...")
	}
}
