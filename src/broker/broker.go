package broker

import (
	"github.com/streadway/amqp"
)

const (
	UnprocessedUrlQueue string = "unprocessed_articles"
	ArticlesQueue       string = "articles"
	RelationsQueue      string = "relations"
)

type Broker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(uri string, queues ...string) (*Broker, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for _, queue := range queues {
		_, err = createQueue(ch, queue)
		if err != nil {
			return nil, err
		}
	}

	return &Broker{conn, ch}, err
}

func createQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // args
	)
}

func (b *Broker) Publish(key string, msg []byte) error {
	return b.ch.Publish(
		"",    // exchange
		key,   // key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
}

func (b *Broker) GetConsumer(queue string) (<-chan amqp.Delivery, error) {
	consumer, err := b.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   //args
	)

	if err != nil {
		return nil, err
	}

	err = b.ch.Qos(1, 0, false)
	return consumer, err
}

func (b *Broker) Ack(tag uint64) error {
	return b.ch.Ack(tag, false)
}
