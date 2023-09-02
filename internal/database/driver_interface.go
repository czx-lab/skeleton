package database

import "gorm.io/gorm"

type DriverInterface interface {
	Instance() gorm.Dialector
}
