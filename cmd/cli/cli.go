package main

import (
	CustomCommand "github.com/czx-lab/skeleton/command"
	_ "github.com/czx-lab/skeleton/internal/bootstrap"
	"github.com/czx-lab/skeleton/internal/command"
)

func main() {
	command.New().AddCommand(
		&CustomCommand.FooCommand{},
	).Execute()
}
