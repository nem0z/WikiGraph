package app

import "github.com/nem0z/WikiGraph/database"

type Config struct {
	BrokerUri      string
	DatabaseConfig *database.Config
}
