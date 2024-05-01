package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	mqbroker "github.com/nem0z/WikiGraph/broker"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("WikiGraph bis")

	err := godotenv.Load()
	handle(err)

	RMQ_URI := os.Getenv("RABBITMQ_URI")
	broker, err := mqbroker.New(RMQ_URI, "testing", "unprocessed_articles", "relations")
	handle(err)

	err = broker.Publish("testing", []byte("Test publishing"))
	handle(err)

	chMsg, err := broker.GetConsumer("testing")
	handle(err)

	for msg := range chMsg {
		fmt.Println(string(msg.Body))
	}

	select {}
}
