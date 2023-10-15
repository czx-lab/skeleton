package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type Command struct {
	root *cobra.Command
}

type Interface interface {
	Command() *cobra.Command
	Flags(*cobra.Command)
}

type CommandInterface interface {
	GlobalFlags()
	RegisterCmds() []Interface
}

func New() *Command {
	root := &Command{
		root: &cobra.Command{
			Use:   "command",
			Short: "skeleton command line.",
			Long: `this command line is an encapsulation of the github.com/spf13/cobra library. 
For the definition of flag, please refer to the official library documentation.`,
			Run: func(cmd *cobra.Command, args []string) {
				Error(cmd, args, errors.New("unrecognized command"))
			},
		},
	}
	return root
}

func (c *Command) Root() *cobra.Command {
	return c.root
}

func (c *Command) AddCommand(cmd CommandInterface) *Command {
	cmd.GlobalFlags()
	commands := cmd.RegisterCmds()
	var cobras []*cobra.Command
	for _, command := range commands {
		cmd := command.Command()
		command.Flags(cmd)
		cobras = append(cobras, cmd)
	}
	c.root.AddCommand(cobras...)
	return c
}

func (c *Command) Execute() {
	_ = c.root.Execute()
}

func ExecuteCommand(name string, subName string, args ...string) (string, error) {
	args = append([]string{subName}, args...)
	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}

func Error(cmd *cobra.Command, args []string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}
