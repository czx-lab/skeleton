package config

import (
	"time"
)

// ConfigInterface 缓存接口
// 应用到redis、cache等cache db
type ConfigInterface interface {
	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
}
