package task

type DemoTask struct {
}

func (d *DemoTask) Rule() string {
	return "* * * * *"
}

func (d *DemoTask) Execute() func() {
	return func() {

	}
}
