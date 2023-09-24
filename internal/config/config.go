package config

import (
	"github.com/spf13/cast"
	"skeleton/internal/constants/config"
	"skeleton/internal/container"
	"sync"
	"time"
)

type Config struct {
	Driver config.DriverInterface
	cache  config.CacheInterface
	mu     *sync.Mutex
}

var _ config.ConfigInterface = (*Config)(nil)

// Options 配置选项
type Options struct {
	// 文件名称
	Filename string

	// 工作目录，项目根目录
	BasePath string

	// 配置文件类型
	Cate string

	// 配置缓存前缀
	CachePrefix string

	// 配置缓存器，可自定义缓存器，比如使用redis
	// 只需要实现CacheInterface接口即可
	Cache config.CacheInterface
}

// New config
func New(config config.DriverInterface, option Options) (provider *Config, err error) {
	if option.Cache == nil {
		option.Cache = container.CreateContainerFactory()
	}
	if option.CachePrefix == "" {
		option.CachePrefix = "config"
	}
	if option.Cate == "" {
		option.Cate = "yaml"
	}
	if d, ok := config.(interface{ Apply(Options) error }); ok {
		if err = d.Apply(option); err != nil {
			return
		}
	}
	if d, ok := config.(interface{ Listen() }); ok {
		d.Listen()
	}
	return &Config{
		Driver: config,
		cache:  option.Cache,
		mu:     new(sync.Mutex),
	}, nil
}

func (c *Config) Cache(key string, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache.Has(key) {
		return true
	}
	return c.cache.Set(key, value)
}

func (c *Config) Get(key string) any {
	if c.cache.Has(key) {
		return c.cache.Get(key)
	}
	val := c.Driver.Get(key)
	c.Cache(key, val)
	return val
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

func (c *Config) GetInt32(key string) int32 {
	return cast.ToInt32(c.Get(key))
}

func (c *Config) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}
