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

func New(cmdFunc func(*cobra.Command) CommandInterface) *Command {
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
	return root.AddCommand(cmdFunc)
}

func (c *Command) Root() *cobra.Command {
	return c.root
}

func (c *Command) AddCommand(cmdFunc func(*cobra.Command) CommandInterface) *Command {
	cmd := cmdFunc(c.root)
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
	if err := c.root.Execute(); err != nil {
		fmt.Printf("Command initialisation or execution failed with an error: %s\n", err)
		os.Exit(1)
	}
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
