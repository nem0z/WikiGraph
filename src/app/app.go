package app

import (
	"log"

	crawlerpkg "github.com/nem0z/WikiGraph/app/crawler"
	"github.com/nem0z/WikiGraph/app/entity"
	"github.com/nem0z/WikiGraph/app/handler"
	brokerpkg "github.com/nem0z/WikiGraph/broker"

	"github.com/nem0z/WikiGraph/database"
)

type App struct {
	crawlers []*crawlerpkg.Crawler
	broker   *brokerpkg.Broker
	db       *database.DB
}

func initCrawlers(broker *brokerpkg.Broker, n int) ([]*crawlerpkg.Crawler, error) {
	crawlers := make([]*crawlerpkg.Crawler, n)

	for i := range crawlers {
		crawler, err := crawlerpkg.New(broker)
		if err != nil {
			return nil, err
		}

		crawlers[i] = crawler
	}

	return crawlers, nil
}

func New(config *Config, nbCrawlers int) (*App, error) {
	broker, err := brokerpkg.New(config.BrokerConfig.Uri(),
		brokerpkg.UnprocessedUrlQueue,
		brokerpkg.ArticlesQueue,
		brokerpkg.RelationsQueue,
	) // TODO : Move queue names in another package

	if err != nil {
		return nil, err
	}

	db, err := database.New(config.DatabaseConfig, publishNewArticleToQueue(broker))
	if err != nil {
		return nil, err
	}

	crawlers, err := initCrawlers(broker, nbCrawlers)
	if err != nil {
		return nil, err
	}

	return &App{
		crawlers: crawlers,
		broker:   broker,
		db:       db,
	}, nil
}

func (app *App) Run() error {
	err := handler.HandleArticles(app.broker, app.db)
	if err != nil {
		return err
	}

	err = handler.HandleRelations(app.broker, app.db)
	if err != nil {
		return err
	}

	for _, crawler := range app.crawlers {
		go crawler.Start()
	}

	return app.broker.Publish(brokerpkg.UnprocessedUrlQueue, []byte("Marseille"))
}

func publishNewArticleToQueue(broker *brokerpkg.Broker) func(article *entity.Article) {
	return func(article *entity.Article) {
		err := broker.Publish(brokerpkg.UnprocessedUrlQueue, []byte(article.Link))
		if err != nil {
			log.Println("Error publishing new article to unprocessed queue")
		}
	}
}
