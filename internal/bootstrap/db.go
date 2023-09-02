package bootstrap

import (
	"github.com/czx-lab/skeleton/internal/database"
	"github.com/czx-lab/skeleton/internal/database/db_log"
	"github.com/czx-lab/skeleton/internal/database/driver"
	"github.com/czx-lab/skeleton/internal/variable"
	"gorm.io/gorm"
	"time"
)

// InitMysql 初始化db实例
// TODO::目前暂时只实现了主从分离，后续加入分库分表
func InitMysql() (*gorm.DB, error) {
	mysqlConfig := variable.Config.Get("Database.Mysql").(map[string]any)
	mysqlMasterDriver := driver.New(driver.WithMysqlDsn(mysqlConfig["write"].(string)))
	dbClass, err := database.New(
		mysqlMasterDriver,
		&gorm.Config{
			Logger: db_log.New(db_log.SetLogger(variable.Log)),
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
