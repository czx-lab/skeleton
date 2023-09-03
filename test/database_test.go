package test

import (
	_ "github.com/czx-lab/skeleton/internal/bootstrap"
	"github.com/czx-lab/skeleton/internal/variable"
	"testing"
)

func TestMysql(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestMysql filed:%v", err)
		}
	}()
	type Result struct {
		ID       int    `gorm:"primaryKey" json:"id"`
		Nickname string `gorm:"column:nickname" json:"nickname"`
		Intro    string `gorm:"column:intro" json:"intro"`
	}

	var result Result
	variable.DB.Raw("SELECT id, nickname, intro FROM user WHERE id = ?", 2).Scan(&result)
	t.Log(result.Intro)
}
