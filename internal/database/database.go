package database

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Database struct {
	Db *gorm.DB

	connMaxLifetime time.Duration
	maxIdleConn     int
	maxOpenConn     int
	connMaxIdleTime time.Duration
}

type Options interface {
	apply(*Database)
}

type OptionFunc func(db *Database)

func New(driver DriverInterface, config *gorm.Config, opts ...Options) (db *Database, err error) {
	gormDb, err := gorm.Open(driver.Instance(), config)
	if err != nil {
		return nil, err
	}
	dbClass := &Database{
		Db: gormDb,
	}
	for _, val := range opts {
		val.apply(dbClass)
	}
	dbClass.hook()
	if err = dbClass.pool(); err != nil {
		return nil, err
	}
	return dbClass, nil
}

func (d *Database) pool() error {
	if rawDb, err := d.Db.DB(); err != nil {
		return err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * d.connMaxIdleTime)
		rawDb.SetConnMaxLifetime(d.connMaxLifetime * time.Second)
		rawDb.SetMaxIdleConns(d.maxIdleConn)
		rawDb.SetMaxOpenConns(d.maxOpenConn)
	}
	return nil
}

func (d *Database) hook() {
	// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
	// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
	_ = d.Db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", MaskNotDataError)
}

func (d *Database) SetReadDb(dbDialectors []gorm.Dialector) error {
	resolverConf := dbresolver.Config{
		Replicas: dbDialectors,
		Policy:   dbresolver.RandomPolicy{},
	}
	if err := d.Db.Use(dbresolver.Register(resolverConf).
		SetConnMaxIdleTime(d.connMaxIdleTime * 30).
		SetConnMaxLifetime(d.connMaxLifetime * time.Second).
		SetMaxIdleConns(d.maxIdleConn).
		SetMaxOpenConns(d.maxOpenConn)); err != nil {
		return err
	}
	return nil
}

func (f OptionFunc) apply(db *Database) {
	f(db)
}

func WithConnMaxLifetime(val time.Duration) Options {
	return OptionFunc(func(db *Database) {
		db.connMaxLifetime = val
	})
}

func WithMaxIdleConn(val int) Options {
	return OptionFunc(func(db *Database) {
		db.maxIdleConn = val
	})
}

func WithMaxOpenConn(val int) Options {
	return OptionFunc(func(db *Database) {
		db.maxOpenConn = val
	})
}

func WithConnMaxIdleTime(val time.Duration) Options {
	return OptionFunc(func(db *Database) {
		db.connMaxIdleTime = val
	})
}
