package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func New(config *Config) (*DB, error) {
	db, err := sql.Open("mysql", config.Uri())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = Init(db, config.InitScriptPath)
	return &DB{db}, err
}
