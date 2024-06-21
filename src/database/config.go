package database

import (
	"fmt"

	"github.com/nem0z/WikiGraph/database/neo4j"
	redispkg "github.com/nem0z/WikiGraph/database/redis"
)

type Config struct {
	User           string
	Pass           string
	Host           string
	DatabaseName   string
	InitScriptPath string
	RedisConfig    *redispkg.Config
	Neo4jConfig    *neo4j.Config
}

func (config *Config) Uri() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.User, config.Pass, config.Host, config.DatabaseName)
}
