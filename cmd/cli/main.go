package main

import (
	CustomCommand "skeleton/app/command"
	_ "skeleton/internal/bootstrap"
	"skeleton/internal/command"
)

func main() {
	cmd := command.New()
	cmd.AddCommand(CustomCommand.NewCommand(cmd.Root())).Execute()
}
