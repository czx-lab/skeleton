package xmysql

import (
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type MySqlConf struct {
	Write string
	Read  []string

	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime int64
	ConnMaxIdleTime int64
}

type XMysql struct {
	db *gorm.DB

	conf MySqlConf
}

func NewMysql(conf MySqlConf) *XMysql {
	instance := &XMysql{conf: conf}

	db, err := gorm.Open(driver(instance.conf.Write))
	if err != nil {
		log.Fatalf("error: mysql init %s", err.Error())
		return nil
	}
	instance.db = db

	instance.hook()
	if err := instance.resolver(); err != nil {
		log.Fatalf("error: mysql init resolver %s", err.Error())
		return nil
	}

	return instance
}

func (m *XMysql) DB() *gorm.DB {
	return m.db
}

func (m *XMysql) hook() {
	// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
	// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
	_ = m.db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", MaskNotDataError)
}

func (m *XMysql) resolver() error {
	if len(m.conf.Read) == 0 {
		return nil
	}
	var dbs []gorm.Dialector
	for _, v := range m.conf.Read {
		dbs = append(dbs, driver(v))
	}
	resolverConf := dbresolver.Config{
		Replicas: dbs,
		Policy:   dbresolver.RandomPolicy{},
	}
	if err := m.db.Use(dbresolver.Register(resolverConf).
		SetConnMaxIdleTime(time.Duration(m.conf.ConnMaxIdleTime) * time.Second).
		SetConnMaxLifetime(time.Duration(m.conf.ConnMaxLifetime) * time.Second).
		SetMaxIdleConns(m.conf.MaxIdleConn).
		SetMaxOpenConns(m.conf.MaxOpenConn)); err != nil {
		return err
	}
	return nil
}

func driver(dsn string) gorm.Dialector {
	return NewMysqlDriver(WithMysqlDsn(dsn), WithMysqlSkipInitializeWithVersion(true)).Instance()
}
