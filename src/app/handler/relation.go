package handler

import (
	"encoding/json"
	"log"

	"github.com/nem0z/WikiGraph/app/entity"
	brokerpkg "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
	"github.com/streadway/amqp"
)

func HandleRelations(broker *brokerpkg.Broker, db *database.DB) error {
	consumer, err := broker.GetConsumer(brokerpkg.RelationsQueue)
	if err != nil {
		return err
	}

	go processRelations(broker, consumer, db)

	return nil
}

func processRelations(broker *brokerpkg.Broker, consumer <-chan amqp.Delivery, db *database.DB) {
	for msg := range consumer {

		var rr *entity.ResolvedRelation
		err := json.Unmarshal(msg.Body, &rr)
		if err != nil {
			log.Println("error unmarshalling relations :", err)
			continue
		}

		err = db.CreateLinks(rr.ParentId, rr.ChildIds...)
		if err != nil {
			log.Println("error creating relations :", err)
			continue
		}

		err = db.ProcessArticle(rr.ParentId)
		if err != nil {
			log.Println("error setting article processed :", err)
			continue
		}

		err = broker.Ack(msg.DeliveryTag)
		if err != nil {
			log.Printf("error acking message %v : %v\n", msg.DeliveryTag, err)
			continue
		}
	}
}
