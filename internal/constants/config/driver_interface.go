package config

import "skeleton/internal/constants/container"

type DriverInterface interface {
	container.ContainerInterface
	Listen()
}
