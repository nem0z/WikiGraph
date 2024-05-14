package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	mqbroker "github.com/nem0z/WikiGraph/broker"
	crawlerpkg "github.com/nem0z/WikiGraph/crawler"
	"github.com/nem0z/WikiGraph/database"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := godotenv.Load()
	handle(err)

	db, err := database.New(nil)
	handle(err)

	rmqUri := os.Getenv("RABBITMQ_URI")
	broker, err := mqbroker.New(rmqUri, crawlerpkg.UnprocessedUrlQueue, crawlerpkg.ArticlesQueue, crawlerpkg.RelationsQueue)
	handle(err)

	err = crawlerpkg.HandleRelations(broker, db)
	handle(err)

	crawler, err := crawlerpkg.New(broker)
	handle(err)

	broker.Publish(crawlerpkg.UnprocessedUrlQueue, []byte("Marseille"))

	go crawler.Work()

	select {}
}
