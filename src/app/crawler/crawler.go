package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nem0z/WikiGraph/app/entity"
	brokerpkg "github.com/nem0z/WikiGraph/broker"
	"github.com/streadway/amqp"
)

type Crawler struct {
	broker   *brokerpkg.Broker
	consumer <-chan amqp.Delivery
	stop     chan bool
	proxy    string
}

func New(broker *brokerpkg.Broker, proxy string) (*Crawler, error) {
	consumer, err := broker.GetConsumer(brokerpkg.UnprocessedUrlQueue)
	if err != nil {
		return nil, err
	}

	chStop := make(chan bool, 1)

	return &Crawler{
		broker:   broker,
		consumer: consumer,
		stop:     chStop,
		proxy:    proxy,
	}, nil
}

func (c *Crawler) Stop() {
	c.stop <- true
}

func (c *Crawler) Start() {
	for {
		select {
		case <-c.stop:
			fmt.Println("Stopping crawler...")
			return
		case msg := <-c.consumer:
			if err := c.process(&msg); err != nil {
				log.Printf("Error processing %v: %v -  acking it\n", msg.DeliveryTag, err)
				c.broker.Ack(msg.DeliveryTag)
			}
		}
	}
}

func (c *Crawler) process(msg *amqp.Delivery) error {
	url := string(msg.Body)
	if url == "" {
		log.Println("Caught empty URL")
		err := c.broker.Ack(msg.DeliveryTag)
		if err != nil {
			log.Printf("error acking message (tag : %v) : %v\n", msg.DeliveryTag, err)
			return err
		}
	}

	scrapper, err := NewScraper(url, c.proxy)
	if err != nil {
		return err
	}

	articles, err := scrapper.GetArticles()
	if err != nil {
		log.Printf("error scrapping articles (url : %v) : %v", url, err)
		return err
	}

	relations := entity.NewRelation(url, articles...)
	bRelations, err := json.Marshal(relations)
	if err != nil {
		log.Printf("error marshalling relations (parent : %v): %v", url, err)
		return err
	}

	err = c.broker.Publish(brokerpkg.ArticlesQueue, bRelations)
	if err != nil {
		log.Printf("error publishing relations: (parent : %v) : %v", url, err)
		return err
	}

	err = c.broker.Ack(msg.DeliveryTag)
	if err != nil {
		log.Printf("error acking message (tag : %v) : %v\n", msg.DeliveryTag, err)
		return err
	}

	return nil
}
