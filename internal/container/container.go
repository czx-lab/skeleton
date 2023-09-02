package container

import (
	"github.com/czx-lab/skeleton/internal/constants/config"
	"strings"
	"sync"
)

var containerMap sync.Map

type Container struct {
}

var _ config.CacheInterface = (*Container)(nil)

func CreateContainerFactory() *Container {
	return &Container{}
}

func (c *Container) Delete(key string) {
	containerMap.Delete(key)
}

func (c *Container) Get(key string) any {
	if value, exist := containerMap.Load(key); exist {
		return value
	}
	return nil
}

func (c *Container) Set(key string, value any) bool {
	var res bool
	if exist := c.Has(key); exist == false {
		containerMap.Store(key, value)
		res = true
	}
	return res
}

func (c *Container) Has(key string) bool {
	_, ok := containerMap.Load(key)
	return ok
}

// FuzzyDelete 按照键的前缀模糊删除容器中注册的内容
func (c *Container) FuzzyDelete(keyPre string) {
	containerMap.Range(func(key, value interface{}) bool {
		if key, ok := key.(string); ok {
			if strings.HasPrefix(key, keyPre) {
				containerMap.Delete(key)
			}
		}
		return true
	})
}
