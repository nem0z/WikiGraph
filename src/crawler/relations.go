package crawler

import (
	"encoding/json"
	"log"
	"time"

	brokerpkg "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
	"github.com/nem0z/WikiGraph/entity"
	"github.com/streadway/amqp"
)

func HandleRelations(broker *brokerpkg.Broker, db *database.DB, n int) error {
	consumer, err := broker.GetConsumer(brokerpkg.RelationsQueue)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		go handleRelations(broker, consumer, db)
	}

	return nil
}

func handleRelations(broker *brokerpkg.Broker, consumer <-chan amqp.Delivery, db *database.DB) {
	for msg := range consumer {
		start := time.Now()

		var relation *entity.Relation
		err := json.Unmarshal(msg.Body, &relation)
		if err != nil {
			log.Println("error unmarshalling relations :", err)
			continue
		}

		err = database.CreateRelations(db, relation)
		if err != nil {
			log.Println("error creating relations :", err)
			continue
		}

		err = broker.Ack(msg.DeliveryTag)
		if err != nil {
			log.Printf("error acking message %v : %v\n", msg.DeliveryTag, err)
			continue
		}

		log.Printf("Creating relations (%v) for : %v", time.Since(start), relation.ParentLink)
	}
}
