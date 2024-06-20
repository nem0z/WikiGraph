package handler

import (
	"log"

	brokerpkg "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
	"github.com/streadway/amqp"
)

func HandleArticles(broker *brokerpkg.Broker, db *database.DB) error {
	consumer, err := broker.GetConsumer(brokerpkg.ArticlesQueue)
	if err != nil {
		return err
	}

	go processArticles(broker, consumer, db)

	return nil
}

func processArticles(broker *brokerpkg.Broker, consumer <-chan amqp.Delivery, db *database.DB) {
	for msg := range consumer {

		url := string(msg.Body)
		if db.Exist(url) {
			continue
		}

		err := broker.Publish(brokerpkg.UnprocessedUrlQueue, msg.Body)
		if err != nil {
			log.Println("error publshing resolved relation")
			continue
		}

		err = broker.Ack(msg.DeliveryTag)
		if err != nil {
			log.Printf("error acking message %v : %v\n", msg.DeliveryTag, err)
			continue
		}
	}
}
