package main

import "github.com/czx-lab/skeleton/internal/command"
import CustomCommand "github.com/czx-lab/skeleton/command"

func main() {
	command.New().AddCommand(
		&CustomCommand.FooCommand{},
	).Execute()
}
