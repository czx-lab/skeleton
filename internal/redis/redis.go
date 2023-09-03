package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Options interface {
	apply(*client)
}

type client struct {
	addr        string
	pwd         string
	db          int
	poolSize    int
	minIdleConn int
	maxIdleConn int
	maxLifetime time.Duration
	maxIdleTime time.Duration
}

func New(opts ...Options) *redis.Client {
	clientClass := &client{}
	for _, val := range opts {
		val.apply(clientClass)
	}
	return redis.NewClient(&redis.Options{
		Addr:            clientClass.addr,
		Password:        clientClass.pwd,
		DB:              clientClass.db,
		PoolSize:        clientClass.poolSize,
		MinIdleConns:    clientClass.minIdleConn,
		MaxIdleConns:    clientClass.maxIdleConn,
		ConnMaxLifetime: clientClass.maxLifetime * time.Second,
		ConnMaxIdleTime: clientClass.maxIdleTime * time.Minute,
	})
}

type Option func(opts *client)

func (f Option) apply(client *client) {
	f(client)
}

func WithAddr(addr string) Option {
	return func(opts *client) {
		opts.addr = addr
	}
}

func WithPwd(pwd string) Option {
	return func(opts *client) {
		opts.pwd = pwd
	}
}

func WithDb(db int) Option {
	return func(opts *client) {
		opts.db = db
	}
}

func WithPoolSize(size int) Option {
	return func(opts *client) {
		opts.poolSize = size
	}
}

func WithMinIdleConn(size int) Option {
	return func(opts *client) {
		opts.minIdleConn = size
	}
}

func WithMaxIdleConn(size int) Option {
	return func(opts *client) {
		opts.maxIdleConn = size
	}
}

func WithMaxLifetime(size time.Duration) Option {
	return func(opts *client) {
		opts.maxLifetime = size
	}
}

func WithMaxIdleTime(size time.Duration) Option {
	return func(opts *client) {
		opts.maxIdleTime = size
	}
}
