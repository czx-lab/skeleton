package command

import (
	AppCommand "skeleton/internal/command"
	"skeleton/internal/variable"

	IDao "skeleton/app/interface/dao"

	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	// https://gorm.io/gen/dao.html#gen-Config
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
	tables = []string{"users"}

	renames = map[string]AppCommand.Rename{
		"users": {
			FileName:  "user",
			ModelName: "userModel",
		},
	}

	// 需要忽略的字段-全局
	ignoreFileds = []string{"updated_at"}

	// Annotation Syntax
	methods = map[string][]any{
		"users": {
			func(IDao.Method) {},
			func(IDao.UserMethod) {},
		},
	}

	// 字段类型转换-全局
	converts = map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return "int8" },
		"timestamp": func(columnType gorm.ColumnType) (dataType string) { return "types.ModelFieldTime" },
		"decimal":   func(columnType gorm.ColumnType) (dataType string) { return "types.Decimal" },
	}
)

type GormGenCommand struct {
	conf         gen.Config
	tables       []string
	ignoreFileds []string
	renames      map[string]AppCommand.Rename
	methods      map[string][]any
	converts     map[string]func(columnType gorm.ColumnType) (dataType string)
}

func NewGormGenCommand() *GormGenCommand {
	return &GormGenCommand{
		conf:         conf,
		tables:       tables,
		renames:      renames,
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
		AppCommand.WithNames(g.renames),
		AppCommand.WithIgnoreFileds(g.ignoreFileds),
		AppCommand.WithDataMap(g.converts),
		AppCommand.WithMethods(g.methods),
	)
}
