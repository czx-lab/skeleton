package config

import (
	"github.com/czx-lab/skeleton/internal/constants/container"
)

// CacheInterface config cache
type CacheInterface interface {
	container.ContainerInterface
	FuzzyDelete(key string)
}
