package broker

import "fmt"

type Config struct {
	User string
	Pass string
	Host string
	Port string
}

func (config *Config) Uri() string {
	return fmt.Sprintf("amqp://%s:%s@%v:%v/",
		config.User, config.Pass, config.Host, config.Port)
}
