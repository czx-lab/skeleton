package test

import (
	"github.com/czx-lab/skeleton/internal/config"
	cinfg "github.com/czx-lab/skeleton/internal/config/driver"
	"github.com/czx-lab/skeleton/internal/database"
	"github.com/czx-lab/skeleton/internal/database/db_log"
	"github.com/czx-lab/skeleton/internal/database/driver"
	"github.com/czx-lab/skeleton/internal/variable"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestMysql(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestMysql filed:%v", err)
		}
	}()
	opt := config.Options{
		BasePath: variable.BasePath,
	}
	provider, _ := config.New(cinfg.New(), opt)
	mysqlConfig := provider.Get("Database.Mysql").(map[string]any)
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
		t.Log(err)
		return
	}
	if read, ok := mysqlConfig["reade"]; ok {
		var dbDialectors []gorm.Dialector
		if read == nil {
			t.Log("success")
			return
		}
		for _, _r := range read.([]any) {
			dbDialectors = append(dbDialectors, driver.New(driver.WithMysqlDsn(_r.(string))).Instance())
		}
		if err := dbClass.SetReadDb(dbDialectors); err != nil {
			t.Log(err)
			return
		}
	}
	type Result struct {
		ID       int    `gorm:"primaryKey" json:"id"`
		Nickname string `gorm:"column:nickname" json:"nickname"`
		Intro    string `gorm:"column:intro" json:"intro"`
	}

	var result Result
	dbClass.Db.Raw("SELECT id, nickname, intro FROM user WHERE id = ?", 2).Scan(&result)
	t.Log(result.Intro)
}
