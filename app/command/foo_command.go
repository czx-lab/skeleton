package command

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

type FooCommand struct {
}

func (*FooCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "hello",
		Short: "A brief description of your command.",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			str, _ := cmd.Flags().GetString("name")
			fmt.Printf("Hello, %s!\n", str)
			t, _ := cmd.Flags().GetBool("time")
			if t {
				fmt.Println("Time:", time.Now().Format("2006-01-02 15:04:05"))
			}
		},
	}
}

func (*FooCommand) Flags(root *cobra.Command) {
	root.Flags().StringP("name", "n", "", "Say hello to someone")
	root.Flags().BoolP("time", "t", false, "Add time info to hello")
}
