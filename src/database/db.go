package database

import (
	"database/sql"

	"github.com/nem0z/WikiGraph/app/entity"
	redispkg "github.com/nem0z/WikiGraph/database/redis"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
	cache           Cache
	onInsertArticle func(article *entity.Article)
}

func New(config *Config, onInsertArticle func(article *entity.Article)) (*DB, error) {
	db, err := sql.Open("mysql", config.Uri())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	cache := redispkg.New(config.RedisConfig)
	return &DB{db, cache, onInsertArticle}, Init(db, config.InitScriptPath)
}
