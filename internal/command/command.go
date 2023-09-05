package command

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

type Command struct {
	root *cobra.Command
}

type Interface interface {
	Command() *cobra.Command
	Flags(*cobra.Command)
}

func New() *Command {
	root := &Command{
		root: &cobra.Command{
			Use:   "command",
			Short: "A brief description of your application.",
			Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications. 
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
			Run: func(cmd *cobra.Command, args []string) {
				Error(cmd, args, errors.New("unrecognized command"))
			},
		},
	}
	return root
}

func (c *Command) AddCommand(commands ...Interface) *Command {
	var cobras []*cobra.Command
	for _, command := range commands {
		cobras = append(cobras, command.Command())
		command.Flags(c.root)
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
