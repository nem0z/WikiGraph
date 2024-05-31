package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nem0z/WikiGraph/app"
	"github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
	"github.com/nem0z/WikiGraph/database/redis"
)

const (
	EnvBrokerHost string = "RABBITMQ_HOST"
	EnvBrokerPort string = "RABBITMQ_PORT"
	EnvBrokerUser string = "RABBITMQ_DEFAULT_USER"
	EnvBrokerPass string = "RABBITMQ_DEFAULT_PASS"

	EnvDatabaseUser    string = "MYSQL_USER"
	EnvDatabasePass    string = "MYSQL_PASSWORD"
	EnvDatabaseHost    string = "MYSQL_HOST"
	EnvDatabaseName    string = "MYSQL_DB"
	InitDatabaseScript string = "init.sql"

	EnvRedisHost string = "REDIS_HOST"
	EnvRedisPort string = "REDIS_PORT"

	DefaultNbCrawlers int = 5

	DotEnvPath string = "../.env"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnv(path string) (*app.Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	brokerHost := os.Getenv(EnvBrokerHost)
	brokerPort := os.Getenv(EnvBrokerPort)
	brokerUser := os.Getenv(EnvBrokerUser)
	brokerPass := os.Getenv(EnvBrokerPass)

	brokerConfig := &broker.Config{
		User: brokerUser,
		Pass: brokerPass,
		Host: brokerHost,
		Port: brokerPort,
	}

	redisHost := os.Getenv(EnvRedisHost)
	redisPort := os.Getenv(EnvRedisPort)

	redisConfig := &redis.Config{Host: redisHost, Port: redisPort}

	dbUser := os.Getenv(EnvDatabaseUser)
	dbPass := os.Getenv(EnvDatabasePass)
	dbHost := os.Getenv(EnvDatabaseHost)
	dbName := os.Getenv(EnvDatabaseName)

	dbConfig := &database.Config{
		User:           dbUser,
		Pass:           dbPass,
		Host:           dbHost,
		DatabaseName:   dbName,
		InitScriptPath: InitDatabaseScript,
		RedisConfig:    redisConfig,
	}

	return &app.Config{
		BrokerConfig:   brokerConfig,
		DatabaseConfig: dbConfig,
	}, nil
}

func main() {
	config, err := loadEnv(DotEnvPath)
	handle(err)

	app, err := app.New(config, DefaultNbCrawlers)
	handle(err)

	app.Run()

	select {}
}
