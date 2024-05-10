package crawler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nem0z/WikiGraph/entity"
)

const baseUrl string = "https://fr.wikipedia.org/wiki/"

type Scraper struct {
	*colly.Collector
}

func NewScraper() *Scraper {
	return &Scraper{colly.NewCollector()}
}

func isValidLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

func FormateUrl(link string) string {
	return strings.TrimPrefix(link, "/wiki/")
}

func (s *Scraper) GetArticles(link string) (articles []*entity.Article, finalError error) {
	s.OnHTML("#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")

		if isValidLink(link) {
			articles = append(articles, entity.NewArticle(FormateUrl(link), title))
		}
	})

	s.OnError(func(r *colly.Response, err error) {
		if err != nil {
			formatedError := fmt.Sprintf("Scrapper on URL : %v failed with response: %v\nError : %v", r.Request.URL, r, err)
			finalError = errors.New(formatedError)
		}
	})

	url := fmt.Sprintf("%v%v", baseUrl, link)
	err := s.Visit(url)
	if finalError != nil {
		finalError = err
	}

	return articles, finalError
}
