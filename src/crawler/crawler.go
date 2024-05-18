package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	mqbroker "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/entity"
	"github.com/streadway/amqp"
)

type Crawler struct {
	scrapper *Scraper
	broker   *mqbroker.Broker
	consumer <-chan amqp.Delivery
	stop     chan bool
}

func New(broker *mqbroker.Broker) (*Crawler, error) {
	chUnprocessedArticles, err := broker.GetConsumer(UnprocessedUrlQueue)
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

	relations := entity.NewRelation(url, articles...)
	bRelations, err := json.Marshal(relations)
	if err != nil {
		log.Printf("error marshalling articles: %v", err)
		return
	}

	err = c.broker.Publish(RelationsQueue, bRelations)
	if err != nil {
		log.Printf("error publishing relations: %v", err)
		return
	}
}
