package command

import (
	"skeleton/app/extend"
	AppCommand "skeleton/internal/command"
	"skeleton/internal/variable"

	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
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

func (c *Command) GlobalFlags() {}

func (c *Command) RegisterCmds() []AppCommand.Interface {
	return []AppCommand.Interface{
		AppCommand.NewServerCommand(),
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
		AppCommand.WithMethods(
			map[string][]any{
				"user": {
					func(extend.Method) {},
					func(extend.UserMethod) {},
				},
			},
		),
		AppCommand.WithDataMap(
			map[string]func(detailType gorm.ColumnType) (dataType string){
				"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int8" },
				"timestamp": func(detailType gorm.ColumnType) (dataType string) { return "extend.LocalTime" },
				"decimal":   func(detailType gorm.ColumnType) (dataType string) { return "extend.Decimal" },
			},
		),
	)
}
