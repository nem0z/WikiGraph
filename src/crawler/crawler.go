package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	mqbroker "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/entity"
	"github.com/streadway/amqp"
)

const (
	UnprocessedArticlesQueue string = "unprocessed_articles"
	ArticlesQueue            string = "articles"
	LinksQueue               string = "links"
)

type Crawler struct {
	scrapper *Scraper
	broker   *mqbroker.Broker
	consumer <-chan amqp.Delivery
	stop     chan bool
}

func New(broker *mqbroker.Broker) (*Crawler, error) {
	chUnprocessedArticles, err := broker.GetConsumer(UnprocessedArticlesQueue)
	if err != nil {
		return nil, err
	}

	chStop := make(chan bool, 1)

	return &Crawler{
		scrapper: NewScraper(),
		broker:   broker,
		consumer: chUnprocessedArticles,
		stop:     chStop,
	}, nil
}

func (c *Crawler) Stop() {
	c.stop <- true
}

func (c *Crawler) Work() {
	for {
		select {
		case <-c.stop:
			fmt.Println("Stopping crawler")
			return
		case msg := <-c.consumer:
			c.work(&msg)
		default:
			continue
		}
	}
}

func (c *Crawler) work(msg *amqp.Delivery) {
	url := string(msg.Body)
	articles, err := c.scrapper.GetArticles(url)
	if err != nil {
		log.Printf("error scrapping articles (%v) : %v", url, err)
		return
	}

	for _, article := range articles {
		bArticle, err := json.Marshal(article)
		if err != nil {
			continue
		}

		err = c.broker.Publish(ArticlesQueue, bArticle)
		if err != nil {
			log.Printf("error publishing article: %v", err)
			continue
		}

		link := entity.NewLink(url, article.Url)
		bLink, err := json.Marshal(link)
		if err != nil {
			continue
		}

		err = c.broker.Publish(LinksQueue, bLink)
		if err != nil {
			log.Printf("error publishing link: %v", err)
			continue
		}
	}
}
