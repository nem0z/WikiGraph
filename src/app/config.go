package app

import "github.com/nem0z/WikiGraph/database"

type Config struct {
	brokerUri      string
	databaseConfig *database.Config
}

func LoadEnvConfig() (*Config, error) {
	return &Config{}, nil // TODO
}
