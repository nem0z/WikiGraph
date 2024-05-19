package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	mqbroker "github.com/nem0z/WikiGraph/broker"
	crawlerpkg "github.com/nem0z/WikiGraph/crawler"
	"github.com/nem0z/WikiGraph/database"
	"github.com/nem0z/WikiGraph/queue"
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
	broker, err := mqbroker.New(rmqUri, mqbroker.UnprocessedUrlQueue, mqbroker.ArticlesQueue, mqbroker.RelationsQueue)
	handle(err)

	q := queue.New(broker, db)
	err = q.Fill()
	handle(err)

	err = crawlerpkg.HandleRelations(broker, db, 1)
	handle(err)

	crawler, err := crawlerpkg.New(broker)
	handle(err)

	go crawler.Work()

	select {}
}
