package queue

import (
	mqbroker "github.com/nem0z/WikiGraph/broker"
	"github.com/nem0z/WikiGraph/database"
)

// TODO : Remove this package to fill the queue when inserting new article

const QueueSize uint = 1000

type Queue struct {
	broker *mqbroker.Broker
	db     *database.DB
}

func New(broker *mqbroker.Broker, db *database.DB) *Queue {
	return &Queue{broker, db}
}

func (q *Queue) Fill() error {
	links, err := database.GetUnprocessedArticleLinks(q.db, QueueSize)
	if err != nil {
		return err
	}

	for _, link := range links {
		q.broker.Publish(mqbroker.UnprocessedUrlQueue, []byte(link))
	}

	return nil
}
