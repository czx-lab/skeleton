package task

import "github.com/czx-lab/skeleton/internal/variable"

func init() {
	if variable.Crontab != nil {
		variable.Crontab.AddFunc(&DemoTask{})
	}
}
