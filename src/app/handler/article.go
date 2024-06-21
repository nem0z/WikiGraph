package handler

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/nem0z/WikiGraph/app/entity"
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

		var relation *entity.Relation
		err := json.Unmarshal(msg.Body, &relation)
		if err != nil {
			log.Println("error unmarshalling relations :", err)
			continue
		}

		for _, article := range relation.Childs {
			err := db.Graph.CreateEdge(relation.ParentLink, article.Link)
			if err != nil {
				log.Println("error creating edges :", err)
				continue
			}

			_, err = db.Cache.Get(article.Link)
			if err == redis.Nil {
				if err := db.Cache.Set(article.Link, true); err != nil {
					log.Println("Cache setting error :", err)
					continue
				}
				db.OnInsertArticle(article)
			} else if err != nil {
				log.Println("Cache error :", err)
			}

		}

		err = broker.Ack(msg.DeliveryTag)
		if err != nil {
			log.Printf("error acking message %v : %v\n", msg.DeliveryTag, err)
			continue
		}
	}
}
