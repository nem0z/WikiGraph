package redis

import "fmt"

type Config struct {
	Host string
	Port string
}

func (config *Config) Uri() string {
	return fmt.Sprintf("%s:%s", config.Host, config.Port)
}
