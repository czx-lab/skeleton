package main

import (
	CustomCommand "github.com/czx-lab/skeleton/app/command"
	_ "github.com/czx-lab/skeleton/internal/bootstrap"
	"github.com/czx-lab/skeleton/internal/command"
)

func main() {
	command.New().AddCommand(
		&CustomCommand.FooCommand{},
	).Execute()
}
