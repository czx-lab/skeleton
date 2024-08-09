package config

import (
	"encoding/json"
	conf "skeleton/internal/config"
)

type Config struct {
	HttpServer struct {
		Port string
		Mode string
	}
}

// MarshalJSON implements config.IAppConfig.
func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

var _ conf.IAppConfig = (*Config)(nil)
