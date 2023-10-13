package command

import (
	"skeleton/app/method"
	AppCommand "skeleton/internal/command"
	"skeleton/internal/variable"

	"github.com/spf13/cobra"
	"gorm.io/gen"
)

type Command struct {
	root *cobra.Command
}

var _ (AppCommand.CommandInterface) = (*Command)(nil)

func NewCommand(root *cobra.Command) *Command {
	return &Command{
		root: root,
	}
}

func (c *Command) GlobalFlags() {
	c.root.PersistentFlags().StringP("foo", "f", "", "foo flag.")
}

func (c *Command) RegisterCmds() []AppCommand.Interface {
	return []AppCommand.Interface{
		&FooCommand{},
		newGenCommand(),
	}
}

func newGenCommand() AppCommand.Interface {
	return AppCommand.NewGenCommand(
		AppCommand.WithConfig(gen.Config{
			OutPath:           "./app/dao",
			OutFile:           "",
			ModelPkgPath:      "./model",
			Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
			FieldNullable:     false,
			FieldCoverable:    false,
			FieldSignable:     false,
			FieldWithIndexTag: false,
			FieldWithTypeTag:  true,
		}),
		AppCommand.WithDB(variable.DB),
		AppCommand.WithTables([]string{"user"}),
		AppCommand.WithIgnoreFileds([]string{"updated_at"}),
		AppCommand.WithMethods(map[string]any{
			"user": func(method.Method) {},
		}),
	)
}