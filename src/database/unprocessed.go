package database

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqbroker "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/crawler"
	"github.com/nem0z/WikiGraph/entity"
	"github.com/streadway/amqp"
)

type MapArticle struct {
	articles map[string]bool
	mu       sync.Mutex
}

func (mapArticle *MapArticle) CountProcessed() int {
	mapArticle.mu.Lock()
	defer mapArticle.mu.Unlock()

	count := 0
	for _, ok := range mapArticle.articles {
		if ok {
			count += 1
		}
	}

	return count
}

func (mapArticle *MapArticle) Display(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		len := len(mapArticle.articles)
		processed := mapArticle.CountProcessed()
		<-ticker.C
		fmt.Printf("Processed %v articles over %v\n", processed, len)
	}
}

func loadArticles() *MapArticle {
	return &MapArticle{
		articles: map[string]bool{
			"France": false,
		},
		mu: sync.Mutex{},
	}
}

func loadUnprocessedArticles(broker *mqbroker.Broker, mapArticle *MapArticle) (finalError error) {
	for url, ok := range mapArticle.articles {
		if !ok {
			err := broker.Publish(crawler.UnprocessedUrlQueue, []byte(url))
			if err != nil {
				finalError = err
			}
		}
	}

	return finalError
}

func routeArticles(broker *mqbroker.Broker, consumer <-chan amqp.Delivery, mapArticles *MapArticle) {
	var article entity.Article
	for msg := range consumer {
		err := json.Unmarshal(msg.Body, &article)
		if err != nil {
			log.Println("error unmarshalling article :", err)
			continue
		}

		mapArticles.mu.Lock()
		if _, ok := mapArticles.articles[article.Url]; !ok {
			mapArticles.articles[article.Url] = true
			broker.Publish(crawler.UnprocessedUrlQueue, []byte(article.Url))
		}
		mapArticles.mu.Unlock()
	}
}

func ArticlesFilterProcess(broker *mqbroker.Broker) (*MapArticle, error) {
	chArticles, err := broker.GetConsumer(crawler.ArticlesQueue)
	if err != nil {
		return nil, err
	}

	mapArticles := loadArticles()
	go loadUnprocessedArticles(broker, mapArticles)
	go routeArticles(broker, chArticles, mapArticles)

	return mapArticles, nil
}
