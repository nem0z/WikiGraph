package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	User           string
	Pass           string
	Host           string
	DatabaseName   string
	InitScriptPath string
}

func DefaultConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := "localhost"
	dbname := os.Getenv("MYSQL_DB")
	initScriptPath := "./init.sql"

	return &Config{user, pass, host, dbname, initScriptPath}, nil
}

func (config *Config) Uri() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.User, config.Pass, config.Host, config.DatabaseName)
}
