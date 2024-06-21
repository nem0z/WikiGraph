package neo4j

import "fmt"

type Config struct {
	Host string
	Port string
	User string
	Pass string
}

func (config *Config) Uri() string {
	return fmt.Sprintf("neo4j://%s:%s", config.Host, config.Port)
}
