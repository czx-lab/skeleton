package command

import (
	"skeleton/app/extend"
	AppCommand "skeleton/internal/command"
	"skeleton/internal/variable"

	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	conf = gen.Config{
		OutPath:           "./app/dao",
		OutFile:           "",
		ModelPkgPath:      "./model",
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	}

	// 生成model的数据表
	tables = []string{"user"}

	// 需要忽略的字段-全局
	ignoreFileds = []string{"updated_at"}

	// Annotation Syntax
	methods = map[string][]any{
		"user": {
			func(extend.Method) {},
			func(extend.UserMethod) {},
		},
	}

	// 字段类型转换-全局
	converts = map[string]func(detailType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int8" },
		"timestamp": func(detailType gorm.ColumnType) (dataType string) { return "extend.LocalTime" },
		"decimal":   func(detailType gorm.ColumnType) (dataType string) { return "extend.Decimal" },
	}
)

type GormGenCommand struct {
	conf         gen.Config
	tables       []string
	ignoreFileds []string
	methods      map[string][]any
	converts     map[string]func(detailType gorm.ColumnType) (dataType string)
}

func NewGormGenCommand() *GormGenCommand {
	return &GormGenCommand{
		conf:         conf,
		tables:       tables,
		ignoreFileds: ignoreFileds,
		methods:      methods,
		converts:     converts,
	}
}

func (g *GormGenCommand) Command() AppCommand.Interface {
	return AppCommand.NewGenCommand(
		AppCommand.WithConfig(g.conf),
		AppCommand.WithDB(variable.DB),
		AppCommand.WithTables(g.tables),
		AppCommand.WithIgnoreFileds(g.ignoreFileds),
		AppCommand.WithDataMap(g.converts),
		AppCommand.WithMethods(g.methods),
	)
}
