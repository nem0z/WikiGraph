package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	brokerpkg "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/entity"
	"github.com/streadway/amqp"
)

type Crawler struct {
	broker   *brokerpkg.Broker
	consumer <-chan amqp.Delivery
	stop     chan bool
}

func New(broker *brokerpkg.Broker) (*Crawler, error) {
	consumer, err := broker.GetConsumer(brokerpkg.UnprocessedUrlQueue)
	if err != nil {
		return nil, err
	}

	chStop := make(chan bool, 1)

	return &Crawler{
		broker:   broker,
		consumer: consumer,
		stop:     chStop,
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
			c.process(&msg)
		default:
			continue
		}
	}
}

func (c *Crawler) process(msg *amqp.Delivery) {
	url := string(msg.Body)
	scrapper := NewScraper(url)

	articles, err := scrapper.GetArticles()
	if err != nil {
		log.Printf("error scrapping articles (url : %v) : %v", url, err)
		return
	}

	relations := entity.NewRelation(url, articles...)
	bRelations, err := json.Marshal(relations)
	if err != nil {
		log.Printf("error marshalling relations (parent : %v): %v", url, err)
		return
	}

	err = c.broker.Publish(brokerpkg.RelationsQueue, bRelations)
	if err != nil {
		log.Printf("error publishing relations: (parent : %v) : %v", url, err)
		return
	}

	err = c.broker.Ack(msg.DeliveryTag)
	if err != nil {
		log.Printf("error acking message (tag : %v) : %v\n", msg.DeliveryTag, err)
		return
	}
}
