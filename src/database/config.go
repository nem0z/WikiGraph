package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	user           string
	pass           string
	host           string
	dbname         string
	initScriptPath string
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
		config.user, config.pass, config.host, config.dbname)
}
