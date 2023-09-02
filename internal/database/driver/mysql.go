package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Options struct {
	mysql.Config
}

type Option func(opts *Options)

type Mysql struct {
	instance gorm.Dialector
	Options  Options
}

func New(opts ...Option) *Mysql {
	mysqlClass := &Mysql{}
	for _, opt := range opts {
		opt(&mysqlClass.Options)
	}
	mysqlClass.instance = mysql.New(mysqlClass.Options.Config)
	return mysqlClass
}

func (m *Mysql) Instance() gorm.Dialector {
	return m.instance
}

// WithMysqlDsn DSN data source name
func WithMysqlDsn(dsn string) Option {
	return func(opts *Options) {
		opts.DSN = dsn
	}
}

// WithMysqlDefaultStringSize string 类型字段的默认长度
func WithMysqlDefaultStringSize(defaultStringSize uint) Option {
	return func(opts *Options) {
		opts.DefaultStringSize = defaultStringSize
	}
}

// WithMysqlDisableDatetimePrecision 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
func WithMysqlDisableDatetimePrecision(disableDatetimePrecision bool) Option {
	return func(opts *Options) {
		opts.DisableDatetimePrecision = disableDatetimePrecision
	}
}

// WithMysqlDontSupportRenameIndex 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
func WithMysqlDontSupportRenameIndex(dontSupportRenameIndex bool) Option {
	return func(opts *Options) {
		opts.DontSupportRenameIndex = dontSupportRenameIndex
	}
}

// WithMysqlDontSupportRenameColumn 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
func WithMysqlDontSupportRenameColumn(dontSupportRenameColumn bool) Option {
	return func(opts *Options) {
		opts.DontSupportRenameColumn = dontSupportRenameColumn
	}
}

// WithMysqlSkipInitializeWithVersion 根据当前 MySQL 版本自动配置
func WithMysqlSkipInitializeWithVersion(skipInitializeWithVersion bool) Option {
	return func(opts *Options) {
		opts.SkipInitializeWithVersion = skipInitializeWithVersion
	}
}
