package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	mqbroker "github.com/nem0z/WikiGraph/broker"
	crawlerpkg "github.com/nem0z/WikiGraph/crawler"
	"github.com/nem0z/WikiGraph/database"
	"github.com/nem0z/WikiGraph/entity"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := godotenv.Load()
	handle(err)

	rmqUri := os.Getenv("RABBITMQ_URI")
	broker, err := mqbroker.New(rmqUri, crawlerpkg.UnprocessedUrlQueue, crawlerpkg.ArticlesQueue, crawlerpkg.LinksQueue)
	handle(err)

	crawler, err := crawlerpkg.New(broker)
	handle(err)

	mapArticles, err := database.ArticlesFilterProcess(broker)
	handle(err)
	go crawler.Work()

	firstArticle := entity.NewArticle("France", "France")
	firstArticleBytes, err := json.Marshal(firstArticle)
	handle(err)

	broker.Publish(crawlerpkg.ArticlesQueue, firstArticleBytes)

	go mapArticles.Display(time.Second)

	select {}
}
