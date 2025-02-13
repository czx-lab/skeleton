package xredis

import (
	"context"
	"fmt"
	"czx/internal/constants"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConf struct {
	Master string
	User   string
	Pass   string
	Addrs  []string
	DB     int

	Pool struct {
		Size        int
		ConnMaxIdle int
		ConnMinIdle int
		MaxLifeTime int
		MaxIdleTime int
	}
}

type XRedis struct {
	conf RedisConf

	instance redis.UniversalClient
}

func New(conf RedisConf) *XRedis {
	xrds := &XRedis{
		conf: conf,
	}
	xrds.instance = xrds.DB()
	return xrds
}

func (x *XRedis) DB() redis.UniversalClient {
	if x.instance == nil {
		x.instance = instance((x.conf))
	}
	return x.instance
}

// Get implements constants.ICache.
func (x *XRedis) Get(key string) (any, error) {
	val, err := x.instance.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		return nil, err
	}
	return val, nil
}

// Has implements constants.ICache.
func (x *XRedis) Has(key string) (bool, error) {
	val, err := x.instance.Exists(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		return false, err
	}
	if err == redis.Nil {
		return false, nil
	}
	return val > 0, nil
}

// Set implements constants.ICache.
func (x *XRedis) Set(key string, value any) (bool, error) {
	if err := x.instance.Set(context.Background(), key, fmt.Sprintf("%v", value), -1).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func instance(conf RedisConf) redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName: conf.Master,
		Username:   conf.User,
		Password:   conf.Pass,
		DB:         conf.DB,
		Addrs:      conf.Addrs,

		PoolFIFO:        true,
		PoolSize:        conf.Pool.Size,
		MinIdleConns:    conf.Pool.ConnMinIdle,
		MaxIdleConns:    conf.Pool.ConnMaxIdle,
		ConnMaxIdleTime: time.Duration(conf.Pool.MaxIdleTime) * time.Second,
		ConnMaxLifetime: time.Duration(conf.Pool.MaxLifeTime) * time.Second,
	})
}

var _ constants.ICache = (*XRedis)(nil)
