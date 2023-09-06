package task

import "fmt"

type DemoTask struct {
}

func (d *DemoTask) Rule() string {
	return "* * * * *"
}

func (d *DemoTask) Execute() func() {
	return func() {
		fmt.Println("demo-task exec...")
	}
}
