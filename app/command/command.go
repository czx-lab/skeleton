package command

import (
	AppCommand "skeleton/internal/command"

	"github.com/spf13/cobra"
)

type Command struct {
	root *cobra.Command
}

var _ (AppCommand.CommandInterface) = (*Command)(nil)

var NewCommand = func(root *cobra.Command) AppCommand.CommandInterface {
	return &Command{
		root: root,
	}
}

// Global flags
func (c *Command) GlobalFlags() {}

func (c *Command) RegisterCmds() []AppCommand.Interface {
	return []AppCommand.Interface{
		AppCommand.NewServerCommand(),
		NewGormGenCommand().Command(),
		&FooCommand{},
	}
}
