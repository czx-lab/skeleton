package bootstrap

import (
	"skeleton/internal/database"
	"skeleton/internal/database/driver"
	dblog "skeleton/internal/database/logger"
	"skeleton/internal/variable"
	"time"

	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

// InitMysql 初始化db实例
// TODO::目前暂时只实现了主从分离，后续加入分库分表
func InitMysql() (*gorm.DB, error) {
	mysqlConfig := variable.Config.Get("Database.Mysql").(map[string]any)
	mysqlMasterDriver := driver.New(driver.WithMysqlDsn(mysqlConfig["write"].(string)))
	slowThreshold := time.Duration(mysqlConfig["slowthreshold"].(int))
	logLevel := gormLog.LogLevel(mysqlConfig["loglevel"].(int))
	dbLog := dblog.New(
		dblog.SetConfig(gormLog.Config{
			SlowThreshold: slowThreshold * time.Second,
			LogLevel:      logLevel,
		}),
		dblog.SetLogger(variable.Log),
	)
	dbClass, err := database.New(
		mysqlMasterDriver,
		&gorm.Config{
			Logger: dbLog,
		},
		database.WithMaxOpenConn(mysqlConfig["maxopenconn"].(int)),
		database.WithMaxIdleConn(mysqlConfig["maxidleconn"].(int)),
		database.WithConnMaxIdleTime(time.Duration(mysqlConfig["connmaxidletime"].(int))),
		database.WithConnMaxLifetime(time.Duration(mysqlConfig["connmaxlifetime"].(int))),
	)
	if err != nil {
		return nil, err
	}
	if read, ok := mysqlConfig["reade"]; ok {
		var dbDialectors []gorm.Dialector
		if read == nil {
			return dbClass.Db, nil
		}
		for _, _r := range read.([]any) {
			dbDialectors = append(dbDialectors, driver.New(driver.WithMysqlDsn(_r.(string))).Instance())
		}
		if err := dbClass.SetReadDb(dbDialectors); err != nil {
			return nil, err
		}
	}
	return dbClass.Db, nil
}
