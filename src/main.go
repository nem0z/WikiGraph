package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nem0z/WikiGraph/app"
	"github.com/nem0z/WikiGraph/database"
)

const (
	EnvBrokerUri       string = "RABBITMQ_URI"
	EnvDatabaseUser    string = "MYSQL_USER"
	EnvDatabasePass    string = "MYSQL_PASSWORD"
	EnvDatabaseHost    string = "MYSQL_HOST"
	EnvDatabaseName    string = "MYSQL_DB"
	InitDatabaseScript string = "init.sql"
	DefaultNbCrawlers  int    = 3
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnv() (*app.Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	brokerUri := os.Getenv(EnvBrokerUri)
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
	config, err := loadEnv()
	handle(err)

	app, err := app.New(config, DefaultNbCrawlers)
	handle(err)

	app.Run()

	select {}
}
