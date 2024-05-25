package app

import (
	"github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
)

type Config struct {
	BrokerConfig   *broker.Config
	DatabaseConfig *database.Config
}
