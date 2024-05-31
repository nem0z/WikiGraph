package handler

import (
	"encoding/json"
	"log"

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

		parentId, err := db.GetIdFromLink(relation.ParentLink)
		if err != nil {
			log.Println("error resolving parent id :", err)
			continue
		}

		childIds, err := db.ResolveArticleIds(relation.Childs...)
		if err != nil {
			log.Println("error resolving childs ids :", err)
			continue
		}

		resolvedRelation := entity.NewResolvedRelation(parentId, childIds...)
		b, err := json.Marshal(resolvedRelation)
		if err != nil {
			log.Println("error marshalling resolved relation")
			continue
		}

		err = broker.Publish(brokerpkg.RelationsQueue, b)
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
