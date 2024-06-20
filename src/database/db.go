package database

import (
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	graph Graph
	cache Cache
}

func New(config *Config) (*DB, error) {
	//graph :=
	//cache := redispkg.New(config.RedisConfig)

	return &DB{}, nil
}

func (db *DB) Exist(key string) bool {
	_, err := db.cache.Get(key)
	return err == nil
}
