package main

import (
	"errors"
	"log"
	"os"
	"strings"

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

	EnvProxyList string = "PROXIES"

	DefaultNbCrawlers int = 10

	DotEnvPath string = "../.env"
)

var proxies = []string{
	"http://23.95.150.145:6114",
	"http://198.23.239.134:6540",
	"http://45.38.107.97:6014",
	"http://207.244.217.165:6712",
	"http://107.172.163.27:6543",
	"http://104.222.161.211:6343",
	"http://64.137.96.74:6641",
	"http://216.10.27.159:6837",
	"http://136.0.207.84:6661",
	"http://142.147.128.93:6593",
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadProxies() ([]string, error) {
	proxiesStr := os.Getenv(EnvProxyList)
	if proxiesStr == "" {
		return nil, errors.New("PROXIES environment variable not set")
	}

	return strings.Split(proxiesStr, ","), nil
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

	proxies, err := loadProxies()
	if err != nil {
		return nil, err
	}

	return &app.Config{
		BrokerConfig:   brokerConfig,
		DatabaseConfig: dbConfig,
		Proxies:        proxies,
	}, nil
}

func main() {
	config, err := loadEnv(DotEnvPath)
	handle(err)

	app, err := app.New(config, DefaultNbCrawlers, proxies)
	handle(err)

	app.Run()

	select {}
}
