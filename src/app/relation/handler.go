package relation

import (
	"encoding/json"
	"log"

	brokerpkg "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
	"github.com/streadway/amqp"
)

func Handle(broker *brokerpkg.Broker, db *database.DB) error {
	consumer, err := broker.GetConsumer(brokerpkg.RelationsQueue)
	if err != nil {
		return err
	}

	go process(broker, consumer, db)

	return nil
}

func process(broker *brokerpkg.Broker, consumer <-chan amqp.Delivery, db *database.DB) {
	for msg := range consumer {

		var relation *Relation
		err := json.Unmarshal(msg.Body, &relation)
		if err != nil {
			log.Println("error unmarshalling relations :", err)
			continue
		}

		err = relation.Create(db)
		if err != nil {
			log.Println("error creating relations :", err)
			continue
		}

		err = broker.Ack(msg.DeliveryTag)
		if err != nil {
			log.Printf("error acking message %v : %v\n", msg.DeliveryTag, err)
			continue
		}
	}
}
