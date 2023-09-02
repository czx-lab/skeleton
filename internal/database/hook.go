package database

import (
	"gorm.io/gorm"
)

type DbHookFunc func(*gorm.DB)

// MaskNotDataError
// 查询无数据时，报错问题（record not found），但是官方认为报错是应该是，我们认为查询无数据，代码一切ok，不应该报错
func MaskNotDataError(db *gorm.DB) {
	db.Statement.RaiseErrorOnNotFound = false
}
