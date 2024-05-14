package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

func (cfg *Config) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.user, cfg.pass, cfg.host, cfg.dbname)
}

type DB struct {
	*sql.DB
}

func New(cfg *Config) (*DB, error) {
	if cfg == nil {
		defaultCfg, err := DefaultConfig()
		if err != nil {
			return nil, err
		}

		cfg = defaultCfg
	}

	db, err := sql.Open("mysql", cfg.String())
	if err != nil {
		return nil, err
	}

	if db.Ping() != nil {
		return nil, errors.New("could not ping database")
	}

	return &DB{db}, Init(db, cfg.initScriptPath)
}
