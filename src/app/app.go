package app

import (
	brokerpkg "github.com/nem0z/WikiGraph/broker"
	crawlerpkg "github.com/nem0z/WikiGraph/crawler"

	"github.com/nem0z/WikiGraph/database"
)

type App struct {
	crawlers []*crawlerpkg.Crawler
	broker   *brokerpkg.Broker
	db       *database.DB
}

func initCrawlers(broker *brokerpkg.Broker, n int) ([]*crawlerpkg.Crawler, error) {
	crawlers := make([]*crawlerpkg.Crawler, 0, n)

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
	broker, err := brokerpkg.New(config.BrokerUri,
		brokerpkg.UnprocessedUrlQueue,
		brokerpkg.ArticlesQueue,
		brokerpkg.RelationsQueue,
	) // TODO : Move queue names in another package

	if err != nil {
		return nil, err
	}

	db, err := database.New(config.DatabaseConfig)
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

func (app *App) Run() {
	for _, crawler := range app.crawlers {
		go crawler.Start()
	}
}
