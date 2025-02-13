package conf

import (
	"log"
	"czx/internal/constants"
	"czx/internal/stores/memo"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	mu    *sync.Mutex
	viper *viper.Viper
	cache constants.ICache

	prefix string
}

var changeTime time.Time

func init() {
	changeTime = time.Now()
}

func New(path string) *Config {
	instance, err := viperInstance(path)
	if err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
		return nil
	}
	conf := &Config{
		viper: instance,
		mu:    new(sync.Mutex),
	}
	defaultConf(conf)

	conf.listen(instance)
	return conf
}

func (c *Config) Load(v any) {
	if err := c.viper.Unmarshal(&v); err != nil {
		log.Fatalf("error: config load %s", err.Error())
	}
}

func (c *Config) Cache(key string, value any) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ok, err := c.cache.Has(key)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return c.cache.Set(key, value)
}

func (c *Config) Get(key string) (any, error) {
	ok, err := c.cache.Has(key)
	if err != nil {
		return nil, err
	}
	if ok {
		return c.cache.Get(key)
	}
	val := c.viper.Get(key)
	c.Cache(key, val)
	return val, nil
}

func (c *Config) listen(instance *viper.Viper) {
	instance.OnConfigChange(func(in fsnotify.Event) {
		if time.Since(changeTime).Seconds() >= 1 {
			if in.Op != fsnotify.Write {
				return
			}
			c.delCache(c.prefix)
			changeTime = time.Now()
		}
	})
}

func (c *Config) delCache(prefix string) {
	method, ok := c.cache.(memo.IMemo)
	if !ok {
		return
	}
	method.FuzzyDel(prefix)
}

func defaultConf(c *Config) {
	if c.cache == nil {
		c.cache = memo.New()
	}
	if len(c.prefix) == 0 {
		c.prefix = "xconf"
	}
}

func viperInstance(path string) (instance *viper.Viper, err error) {
	instance = viper.New()
	instance.SetConfigFile(path)
	if err = instance.ReadInConfig(); err != nil {
		return
	}
	return
}
