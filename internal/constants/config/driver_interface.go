package config

import "github.com/czx-lab/skeleton/internal/constants/container"

type DriverInterface interface {
	container.ContainerInterface
	Listen()
}
