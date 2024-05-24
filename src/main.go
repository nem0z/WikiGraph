package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nem0z/WikiGraph/app"
	"github.com/nem0z/WikiGraph/database"
)

const (
	EnvBrokerUser      string = "RABBITMQ_DEFAULT_USER"
	EnvBrokerPass      string = "RABBITMQ_DEFAULT_PASS"
	EnvDatabaseUser    string = "MYSQL_USER"
	EnvDatabasePass    string = "MYSQL_PASSWORD"
	EnvDatabaseHost    string = "MYSQL_HOST"
	EnvDatabaseName    string = "MYSQL_DB"
	InitDatabaseScript string = "init.sql"
	DefaultNbCrawlers  int    = 3

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

	brokerUser := os.Getenv(EnvBrokerUser)
	brokerPass := os.Getenv(EnvBrokerPass)
	brokerUri := fmt.Sprintf("amqp://%s:%s@localhost:5672/", brokerUser, brokerPass)

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
	}

	return &app.Config{
		BrokerUri:      brokerUri,
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
